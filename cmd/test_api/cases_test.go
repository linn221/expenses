package main

// Case is one HTTP call. Edit Cases to change what the suite runs;
// the runner in api_test.go stays as it is.
//
// Path, Body, and Want values can contain {{name}} placeholders. A placeholder
// is filled from a variable saved by an earlier case.
//
// Auth sends the saved "session" variable as the session_id header. The login
// case must save the response field session_id under the name session.
//
// Set Skip to leave a case in the list without calling it.
// Set Only on a case to run the sequence from the start through that case
// and ignore everything after it. When several cases set Only, the run stops
// after the last one of them.
type Case struct {
	Name       string
	Method     string
	Path       string
	Body       string
	Auth       bool
	WantStatus int
	WantBody   string
	Save       map[string]string
	Want       map[string]string
	Skip       bool
	Only       bool
}

// Cases walks every route once, in an order that satisfies foreign keys.
// Change this slice to point the same runner at a different scenario.
var Cases = []Case{
	{
		Name:       "register",
		Method:     "POST",
		Path:       "/register",
		Body:       `{"name":"user","email":"user@example.com","password":"secret44"}`,
		WantStatus: 200,
		Want:       map[string]string{"name": "user", "email": "user@example.com"},
	},
	{
		Name:       "categories require a session",
		Method:     "GET",
		Path:       "/categories",
		WantStatus: 401,
	},
	{
		Name:       "login",
		Method:     "POST",
		Path:       "/login",
		Body:       `{"email":"user@example.com","password":"secret44"}`,
		WantStatus: 200,
		Save:       map[string]string{"session": "session_id"},
	},
	{
		Name:       "root",
		Method:     "GET",
		Path:       "/",
		WantStatus: 200,
		WantBody:   "hello world",
	},

	{
		Name:       "create category",
		Method:     "POST",
		Path:       "/categories",
		Auth:       true,
		Body:       `{"name":"groceries"}`,
		WantStatus: 201,
		Save:       map[string]string{"categoryId": "id"},
		Want:       map[string]string{"name": "groceries"},
	},
	{
		Name:       "get category",
		Method:     "GET",
		Path:       "/categories/{{categoryId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"id": "{{categoryId}}", "name": "groceries"},
	},
	{
		Name:       "update category",
		Method:     "PUT",
		Path:       "/categories/{{categoryId}}",
		Auth:       true,
		Body:       `{"name":"food"}`,
		WantStatus: 200,
	},
	{
		Name:       "get category after update",
		Method:     "GET",
		Path:       "/categories/{{categoryId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"name": "food"},
	},
	{
		Name:       "list categories",
		Method:     "GET",
		Path:       "/categories?name=food",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"results.0.id": "{{categoryId}}", "results.0.name": "food"},
	},
	{
		Name:       "paginate categories",
		Method:     "GET",
		Path:       "/categories/paginate?name=food&page=1&order=id&direction=asc",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"page": "1", "per_page": "20", "total": "1", "data.results.0.name": "food"},
	},
	{
		Name:       "deactivate category",
		Method:     "POST",
		Path:       "/is-active/categories/{{categoryId}}?is_active=0",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "activate category",
		Method:     "POST",
		Path:       "/is-active/categories/{{categoryId}}?is_active=1",
		Auth:       true,
		WantStatus: 204,
	},

	{
		Name:       "list wallets",
		Method:     "GET",
		Path:       "/wallets",
		Auth:       true,
		WantStatus: 200,
		Save:       map[string]string{"walletId": "0.id"},
		Want:       map[string]string{"0.name": "Wallet", "0.currency": "MMK"},
	},
	{
		Name:       "create wallet",
		Method:     "POST",
		Path:       "/wallets",
		Auth:       true,
		Body:       `{"name":"cash","currency":"USD"}`,
		WantStatus: 201,
		Save:       map[string]string{"extraWalletId": "id"},
		Want:       map[string]string{"name": "cash", "currency": "USD"},
	},
	{
		Name:       "get wallet",
		Method:     "GET",
		Path:       "/wallets/{{extraWalletId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"id": "{{extraWalletId}}", "name": "cash", "currency": "USD"},
	},
	{
		Name:       "update wallet",
		Method:     "PUT",
		Path:       "/wallets/{{extraWalletId}}",
		Auth:       true,
		Body:       `{"name":"card","currency":"EUR"}`,
		WantStatus: 200,
	},
	{
		Name:       "get wallet after update",
		Method:     "GET",
		Path:       "/wallets/{{extraWalletId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"name": "card", "currency": "EUR"},
	},
	{
		Name:       "deactivate wallet",
		Method:     "POST",
		Path:       "/is-active/wallets/{{extraWalletId}}?is_active=0",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "activate wallet",
		Method:     "POST",
		Path:       "/is-active/wallets/{{extraWalletId}}?is_active=1",
		Auth:       true,
		WantStatus: 204,
	},

	{
		Name:       "create item",
		Method:     "POST",
		Path:       "/items",
		Auth:       true,
		Body:       `{"name":"coffee","code":"coffee","price":"400","category_id":{{categoryId}}}`,
		WantStatus: 201,
		Save:       map[string]string{"itemId": "id"},
		Want:       map[string]string{"name": "coffee", "code": "coffee", "price": "400"},
	},
	{
		Name:       "get item",
		Method:     "GET",
		Path:       "/items/{{itemId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"id": "{{itemId}}", "name": "coffee", "category": "food"},
	},
	{
		Name:       "update item",
		Method:     "PUT",
		Path:       "/items/{{itemId}}",
		Auth:       true,
		Body:       `{"name":"latte","code":"latte","price":"450","category_id":{{categoryId}}}`,
		WantStatus: 200,
		Want:       map[string]string{"name": "latte", "price": "450", "category": "food"},
	},
	{
		Name:       "list items",
		Method:     "GET",
		Path:       "/items?search=latte&category_id={{categoryId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"items.0.id": "{{itemId}}", "items.0.name": "latte", "items.0.category": "food"},
	},
	{
		Name:       "paginate items",
		Method:     "GET",
		Path:       "/items/paginate?search=latte&page=1&order=id&direction=asc",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"page": "1", "total": "1", "data.items.0.name": "latte"},
	},
	{
		Name:       "deactivate item",
		Method:     "POST",
		Path:       "/is-active/items/{{itemId}}?is_active=0",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "activate item",
		Method:     "POST",
		Path:       "/is-active/items/{{itemId}}?is_active=1",
		Auth:       true,
		WantStatus: 204,
	},

	{
		Name:       "create income",
		Method:     "POST",
		Path:       "/incomes",
		Auth:       true,
		Body:       `{"title":"salary","amount":"500000","income_date":"2026-10-28","wallet_id":{{walletId}}}`,
		WantStatus: 201,
		Save:       map[string]string{"incomeId": "id"},
		Want:       map[string]string{"title": "salary", "amount": "500000", "wallet.id": "{{walletId}}", "wallet.name": "Wallet"},
	},
	{
		Name:       "get income",
		Method:     "GET",
		Path:       "/incomes/{{incomeId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"id": "{{incomeId}}", "title": "salary", "amount": "500000"},
	},
	{
		Name:       "update income",
		Method:     "PUT",
		Path:       "/incomes/{{incomeId}}",
		Auth:       true,
		Body:       `{"title":"bonus","amount":"400000","income_date":"2026-10-28","wallet_id":{{walletId}}}`,
		WantStatus: 200,
	},
	{
		Name:       "get income after update",
		Method:     "GET",
		Path:       "/incomes/{{incomeId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"title": "bonus", "amount": "400000"},
	},
	{
		Name:       "list incomes",
		Method:     "GET",
		Path:       "/incomes/list/{{walletId}}?search=bonus&income_date=2026-10-28",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"incomes.0.id": "{{incomeId}}", "incomes.0.title": "bonus", "incomes.0.amount": "400000"},
	},
	{
		Name:       "paginate incomes",
		Method:     "GET",
		Path:       "/incomes/paginate/{{walletId}}?page=1&order=id&direction=asc",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"page": "1", "total": "1", "data.incomes.0.title": "bonus"},
	},

	{
		Name:       "create expense",
		Method:     "POST",
		Path:       "/expenses",
		Auth:       true,
		Body:       `{"expense_date":"2026-10-28","note":"morning coffee","WalletId":{{walletId}},"item_id":{{itemId}},"qty":1,"item_price":"500","UpdateItemPrice":false}`,
		WantStatus: 201,
		Save:       map[string]string{"expenseId": "id"},
		Want:       map[string]string{"note": "morning coffee", "amount": "500", "expense_date": "2026-10-28"},
	},
	{
		Name:       "get expense",
		Method:     "GET",
		Path:       "/expenses/{{expenseId}}",
		Auth:       true,
		WantStatus: 200,
		Want: map[string]string{
			"id":            "{{expenseId}}",
			"note":          "morning coffee",
			"amount":        "500",
			"wallet.id":     "{{walletId}}",
			"wallet.name":   "Wallet",
			"item.id":       "{{itemId}}",
			"item.name":     "latte",
			"item.category": "food",
		},
	},
	{
		Name:       "update expense",
		Method:     "PUT",
		Path:       "/expenses/{{expenseId}}",
		Auth:       true,
		Body:       `{"expense_date":"2026-10-29","note":"evening coffee","WalletId":{{walletId}},"item_id":{{itemId}},"qty":2,"item_price":"500","UpdateItemPrice":false}`,
		WantStatus: 200,
	},
	{
		Name:       "get expense after update",
		Method:     "GET",
		Path:       "/expenses/{{expenseId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"note": "evening coffee", "amount": "1000", "expense_date": "2026-10-29"},
	},
	{
		Name:       "list expenses",
		Method:     "GET",
		Path:       "/expenses/list/{{walletId}}?search=evening&expense_date=2026-10-29",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"expenses.0.id": "{{expenseId}}", "expenses.0.note": "evening coffee", "expenses.0.amount": "1000"},
	},
	{
		Name:       "paginate expenses",
		Method:     "GET",
		Path:       "/expenses/paginate/{{walletId}}?page=1&order=id&direction=desc",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"page": "1", "total": "1", "data.expenses.0.note": "evening coffee"},
	},

	{
		Name:   "create cost perspective",
		Method: "POST",
		Path:   "/cost-perspectives",
		Auth:   true,
		Body: `{
			"name": "markup",
			"value": "15",
			"type": "percentage",
			"applications": [
				{"money_insight_category": ["expense", "daily_expense"], "starting_value": "1000"}
			]
		}`,
		WantStatus: 201,
		Save:       map[string]string{"costPerspectiveId": "id"},
		Want: map[string]string{
			"name":                                  "markup",
			"value":                                 "15",
			"type":                                  "percentage",
			"applications.0.money_insight_category": "expense",
			"applications.0.starting_value":         "1000",
			"applications.1.money_insight_category": "daily_expense",
		},
	},
	{
		Name:       "get cost perspective",
		Method:     "GET",
		Path:       "/cost-perspectives/{{costPerspectiveId}}",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"id": "{{costPerspectiveId}}", "name": "markup", "type": "percentage"},
	},
	{
		Name:   "update cost perspective",
		Method: "PUT",
		Path:   "/cost-perspectives/{{costPerspectiveId}}",
		Auth:   true,
		Body: `{
			"name": "fees",
			"value": "20",
			"type": "ratio",
			"applications": [
				{"money_insight_category": ["income"], "starting_value": "2000"}
			]
		}`,
		WantStatus: 200,
	},
	{
		Name:       "get cost perspective after update",
		Method:     "GET",
		Path:       "/cost-perspectives/{{costPerspectiveId}}",
		Auth:       true,
		WantStatus: 200,
		Want: map[string]string{
			"name":                                  "fees",
			"value":                                 "20",
			"type":                                  "ratio",
			"applications.0.money_insight_category": "income",
			"applications.0.starting_value":         "2000",
		},
	},
	{
		Name:       "list cost perspectives",
		Method:     "GET",
		Path:       "/cost-perspectives?name=fees&type=ratio",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"results.0.id": "{{costPerspectiveId}}", "results.0.name": "fees"},
	},
	{
		Name:       "paginate cost perspectives",
		Method:     "GET",
		Path:       "/cost-perspectives/paginate?name=fees&page=1&order=id&direction=asc",
		Auth:       true,
		WantStatus: 200,
		Want:       map[string]string{"page": "1", "total": "1", "data.results.0.name": "fees"},
	},
	{
		Name:       "deactivate cost perspective",
		Method:     "POST",
		Path:       "/is-active/cost-perspectives/{{costPerspectiveId}}?is_active=0",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "activate cost perspective",
		Method:     "POST",
		Path:       "/is-active/cost-perspectives/{{costPerspectiveId}}?is_active=1",
		Auth:       true,
		WantStatus: 204,
	},

	{
		Name:       "delete expense",
		Method:     "DELETE",
		Path:       "/expenses/{{expenseId}}",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "delete income",
		Method:     "DELETE",
		Path:       "/incomes/{{incomeId}}",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "delete item",
		Method:     "DELETE",
		Path:       "/items/{{itemId}}",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "delete category",
		Method:     "DELETE",
		Path:       "/categories/{{categoryId}}",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "delete extra wallet",
		Method:     "DELETE",
		Path:       "/wallets/{{extraWalletId}}",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "delete cost perspective",
		Method:     "DELETE",
		Path:       "/cost-perspectives/{{costPerspectiveId}}",
		Auth:       true,
		WantStatus: 204,
	},
	{
		Name:       "delete default wallet",
		Method:     "DELETE",
		Path:       "/wallets/{{walletId}}",
		Auth:       true,
		WantStatus: 204,
	},
}
