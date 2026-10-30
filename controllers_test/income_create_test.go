package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/controllers"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/services"
	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateIncomeRespondsCreated(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Income{}); err != nil {
		t.Fatal(err)
	}

	user := models.User{Name: "linn", Email: "linn@example.com", Password: "secret"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	wallet := models.Wallet{
		Name:        "cash",
		Currency:    "MMK",
		Current:     decimal.NewFromInt(1000),
		HasUserId:   models.HasUserId{UserId: user.Id},
		HasIsActive: models.HasIsActive{IsActive: true},
	}
	if err := db.Create(&wallet).Error; err != nil {
		t.Fatal(err)
	}

	container := &app.Container{DB: db}
	provider := services.New(container)
	income := controllers.New(container, provider).IncomeController

	body := bytes.NewBufferString(`{
		"title": "salary",
		"amount": "500000",
		"income_date": "2026-10-28",
		"wallet_id": ` + strconv.Itoa(wallet.Id) + `
	}`)
	req := httptest.NewRequest(http.MethodPost, "/incomes", body)
	req = req.WithContext(app.SetUserId(context.Background(), user.Id))
	rec := httptest.NewRecorder()

	income.Create().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Id     int    `json:"id"`
		Title  string `json:"title"`
		Amount string `json:"amount"`
		Wallet struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"wallet"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Title != "salary" || resp.Wallet.Name != "cash" || resp.Wallet.Id != wallet.Id {
		t.Fatalf("response = %+v", resp)
	}

	var stored models.Wallet
	if err := db.First(&stored, wallet.Id).Error; err != nil {
		t.Fatal(err)
	}
	if !stored.Current.Equal(decimal.NewFromInt(501000)) {
		t.Fatalf("wallet current = %s, want 501000", stored.Current)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/incomes/"+strconv.Itoa(resp.Id), nil)
	getReq.SetPathValue("id", strconv.Itoa(resp.Id))
	getReq = getReq.WithContext(app.SetUserId(context.Background(), user.Id))
	getRec := httptest.NewRecorder()
	income.Get().ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("get status = %d, body = %s", getRec.Code, getRec.Body.String())
	}
}
