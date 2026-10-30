package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/linn221/expenses/routing"
)

func TestAPI(t *testing.T) {
	server := startServer(t)
	client := &http.Client{}
	vars := map[string]string{}

	for _, tc := range casesToRun(Cases) {
		if t.Failed() {
			return
		}
		t.Run(tc.Name, func(t *testing.T) {
			runCase(t, client, server.URL, vars, tc)
		})
	}
}

func startServer(t *testing.T) *httptest.Server {
	t.Helper()
	t.Setenv("SQLITE_PATH", filepath.Join(t.TempDir(), "expenses.db")+"?_foreign_keys=on")

	db := newContainer()
	sqlDB, err := db.DB.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sqlDB.Close() })

	server := httptest.NewServer(routing.Handler(db))
	t.Cleanup(server.Close)
	return server
}

func casesToRun(all []Case) []Case {
	lastOnly := -1
	for i, c := range all {
		if c.Only {
			lastOnly = i
		}
	}
	if lastOnly >= 0 {
		all = all[:lastOnly+1]
	}
	selected := make([]Case, 0, len(all))
	for _, c := range all {
		if !c.Skip {
			selected = append(selected, c)
		}
	}
	return selected
}

var placeholder = regexp.MustCompile(`\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}`)

func substitute(t *testing.T, value string, vars map[string]string) string {
	t.Helper()
	missing := ""
	replaced := placeholder.ReplaceAllStringFunc(value, func(match string) string {
		name := placeholder.FindStringSubmatch(match)[1]
		saved, ok := vars[name]
		if !ok {
			missing = match
			return match
		}
		return saved
	})
	if missing != "" {
		t.Fatalf("unknown placeholder %s", missing)
	}
	return replaced
}

func runCase(t *testing.T, client *http.Client, base string, vars map[string]string, tc Case) {
	t.Helper()

	var body io.Reader
	if tc.Body != "" {
		body = strings.NewReader(substitute(t, tc.Body, vars))
	}
	req, err := http.NewRequest(tc.Method, base+substitute(t, tc.Path, vars), body)
	if err != nil {
		t.Fatal(err)
	}
	if tc.Body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if tc.Auth {
		session, ok := vars["session"]
		if !ok || session == "" {
			t.Fatal("Auth is set, and no session has been saved; the login case must save session_id as session")
		}
		req.Header.Set("session_id", session)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != tc.WantStatus {
		t.Fatalf("status = %d, want %d, body = %s", resp.StatusCode, tc.WantStatus, snippet(payload))
	}
	if tc.WantBody != "" && string(payload) != tc.WantBody {
		t.Fatalf("body = %q, want %q", payload, tc.WantBody)
	}
	if len(tc.Save) == 0 && len(tc.Want) == 0 {
		return
	}

	decoded := decodeJSON(t, payload)
	for name, path := range tc.Save {
		got, ok := jsonValue(decoded, path)
		if !ok {
			t.Fatalf("save %s: path %s not found in %s", name, path, snippet(payload))
		}
		vars[name] = got
	}
	for path, want := range tc.Want {
		got, ok := jsonValue(decoded, path)
		if !ok {
			t.Fatalf("want %s: path not found in %s", path, snippet(payload))
		}
		want = substitute(t, want, vars)
		if got != want {
			t.Fatalf("path %s = %s, want %s", path, got, want)
		}
	}
}

func decodeJSON(t *testing.T, payload []byte) any {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var decoded any
	if err := dec.Decode(&decoded); err != nil {
		t.Fatalf("response is not json: %s", snippet(payload))
	}
	return decoded
}

func jsonValue(node any, path string) (string, bool) {
	current := node
	if path != "" {
		for _, part := range strings.Split(path, ".") {
			switch typed := current.(type) {
			case map[string]any:
				next, ok := typed[part]
				if !ok {
					return "", false
				}
				current = next
			case []any:
				index, err := strconv.Atoi(part)
				if err != nil || index < 0 || index >= len(typed) {
					return "", false
				}
				current = typed[index]
			default:
				return "", false
			}
		}
	}
	switch typed := current.(type) {
	case string:
		return typed, true
	case json.Number:
		return typed.String(), true
	case bool:
		return strconv.FormatBool(typed), true
	default:
		return "", false
	}
}

func snippet(payload []byte) string {
	text := strings.TrimSpace(string(payload))
	if len(text) > 500 {
		return text[:500] + "..."
	}
	return text
}
