### Authentication

Resource routes require a session. Send the id from login as a `session_id` cookie or header. The session lasts 7 days.

```http
session_id: 7c9e6679-7425-40de-944b-e07fc1f90ae7
```

Register, login, and `GET /` are public. Register and login allow 50 calls per IP every 2 minutes. Amounts are JSON strings. Dates are `YYYY-MM-DD`.

# Documentation

- [1. Authentication](#1-authentication)
   * [1.1. register](#11-register)
   * [1.2. login](#12-login)
   * [1.3. hello](#13-hello)
- [2. categories](#2-categories)
   * [2.1. create category](#21-create-category)
   * [2.2. show category](#22-show-category)
   * [2.3. update category](#23-update-category)
   * [2.4. delete category](#24-delete-category)
   * [2.5. list categories](#25-list-categories)
   * [2.6. paginate categories](#26-paginate-categories)
   * [2.7. toggle category](#27-toggle-category)
- [3. wallets](#3-wallets)
   * [3.1. create wallet](#31-create-wallet)
   * [3.2. show wallet](#32-show-wallet)
   * [3.3. update wallet](#33-update-wallet)
   * [3.4. delete wallet](#34-delete-wallet)
   * [3.5. list wallets](#35-list-wallets)
   * [3.6. toggle wallet](#36-toggle-wallet)
- [4. items](#4-items)
   * [4.1. create item](#41-create-item)
   * [4.2. show item](#42-show-item)
   * [4.3. update item](#43-update-item)
   * [4.4. delete item](#44-delete-item)
   * [4.5. list items](#45-list-items)
   * [4.6. paginate items](#46-paginate-items)
   * [4.7. toggle item](#47-toggle-item)
- [5. incomes](#5-incomes)
   * [5.1. create income](#51-create-income)
   * [5.2. show income](#52-show-income)
   * [5.3. update income](#53-update-income)
   * [5.4. delete income](#54-delete-income)
   * [5.5. list incomes](#55-list-incomes)
   * [5.6. paginate incomes](#56-paginate-incomes)
- [6. expenses](#6-expenses)
   * [6.1. create expense](#61-create-expense)
   * [6.2. show expense](#62-show-expense)
   * [6.3. update expense](#63-update-expense)
   * [6.4. delete expense](#64-delete-expense)
   * [6.5. list expenses](#65-list-expenses)
   * [6.6. paginate expenses](#66-paginate-expenses)
- [7. cost-perspectives](#7-cost-perspectives)
   * [7.1. create cost perspective](#71-create-cost-perspective)
   * [7.2. show cost perspective](#72-show-cost-perspective)
   * [7.3. update cost perspective](#73-update-cost-perspective)
   * [7.4. delete cost perspective](#74-delete-cost-perspective)
   * [7.5. list cost perspectives](#75-list-cost-perspectives)
   * [7.6. paginate cost perspectives](#76-paginate-cost-perspectives)
   * [7.7. toggle cost perspective](#77-toggle-cost-perspective)

# 1. Authentication
## 1.1. register

Note:
also creates a wallet named Wallet in MMK. Name is 3–255 characters, password is 8–255, email must be unique.

**POST**
```http
http://localhost:8080/register
```
[*raw*]

```json
{
    "name": "user",
    "email": "user@example.com",
    "password": "secret44"
}
```


## 1.2. login

Note:
returns `session_id` and sets the cookie.

**POST**
```http
http://localhost:8080/login
```
[*raw*]

```json
{
    "email": "user@example.com",
    "password": "secret44"
}
```


## 1.3. hello


**GET**
```http
http://localhost:8080/
```

----------------------------
# 2. categories
## 2.1. create category

Note:
name is required, 3–255 characters, and unique for the user.

**POST**
```http
http://localhost:8080/categories
```
[*raw*]

```json
{
    "name": "groceries"
}
```


## 2.2. show category


**GET**
```http
http://localhost:8080/categories/1
```

## 2.3. update category


**PUT**
```http
http://localhost:8080/categories/1
```
[*raw*]

```json
{
    "name": "food"
}
```


## 2.4. delete category

Note:
rejected while an item uses the category.

**DELETE**
```http
http://localhost:8080/categories/1
```

## 2.5. list categories

Note:
optional `name` and `is_active` (`1` is true).

**GET**
```http
http://localhost:8080/categories?name=food
```

## 2.6. paginate categories

Note:
`page` starts at 1. `order` is a column name. `direction` is `asc` or `desc`. Page size is 20.

**GET**
```http
http://localhost:8080/categories/paginate?name=food&page=1&order=id&direction=asc
```

## 2.7. toggle category

Note:
`is_active=1` turns it on. Any other value turns it off. `204` when it changes, `304` when it does not. Rejected while an item uses it.

**POST**
```http
http://localhost:8080/is-active/categories/1?is_active=0
```

----------------------------
# 3. wallets
## 3.1. create wallet

Note:
name and currency are each unique for the user.

**POST**
```http
http://localhost:8080/wallets
```
[*raw*]

```json
{
    "name": "cash",
    "currency": "USD"
}
```


## 3.2. show wallet


**GET**
```http
http://localhost:8080/wallets/2
```

## 3.3. update wallet

Note:
does not change the balance.

**PUT**
```http
http://localhost:8080/wallets/2
```
[*raw*]

```json
{
    "name": "card",
    "currency": "EUR"
}
```


## 3.4. delete wallet

Note:
rejected while an income or expense uses it.

**DELETE**
```http
http://localhost:8080/wallets/2
```

## 3.5. list wallets


**GET**
```http
http://localhost:8080/wallets
```

## 3.6. toggle wallet

Note:
`is_active=1` turns it on. Cannot turn off a wallet that has an expense.

**POST**
```http
http://localhost:8080/is-active/wallets/2?is_active=0
```

----------------------------
# 4. items
## 4.1. create item

Note:
price must be greater than zero. `category_id` must belong to the user. `code` is optional and unique; if omitted it is the lowercased name, cut to 6 characters.

**POST**
```http
http://localhost:8080/items
```
[*raw*]

```json
{
    "name": "coffee",
    "code": "coffee",
    "price": "400",
    "category_id": 1
}
```


## 4.2. show item


**GET**
```http
http://localhost:8080/items/1
```

## 4.3. update item


**PUT**
```http
http://localhost:8080/items/1
```
[*raw*]

```json
{
    "name": "latte",
    "code": "latte",
    "price": "450",
    "category_id": 1
}
```


## 4.4. delete item


**DELETE**
```http
http://localhost:8080/items/1
```

## 4.5. list items

Note:
optional `search` (name or code) and `category_id`.

**GET**
```http
http://localhost:8080/items?search=latte&category_id=1
```

## 4.6. paginate items


**GET**
```http
http://localhost:8080/items/paginate?search=latte&page=1&order=id&direction=asc
```

## 4.7. toggle item


**POST**
```http
http://localhost:8080/is-active/items/1?is_active=0
```

----------------------------
# 5. incomes
## 5.1. create income

Note:
amount must be greater than zero. Adds the amount to the wallet balance.

**POST**
```http
http://localhost:8080/incomes
```
[*raw*]

```json
{
    "title": "salary",
    "amount": "500000",
    "income_date": "2026-10-28",
    "wallet_id": 1
}
```


## 5.2. show income


**GET**
```http
http://localhost:8080/incomes/1
```

## 5.3. update income

Note:
adjusts the wallet balance by the difference. Does not move the income to another wallet.

**PUT**
```http
http://localhost:8080/incomes/1
```
[*raw*]

```json
{
    "title": "bonus",
    "amount": "400000",
    "income_date": "2026-10-28",
    "wallet_id": 1
}
```


## 5.4. delete income

Note:
subtracts the amount from the wallet.

**DELETE**
```http
http://localhost:8080/incomes/1
```

## 5.5. list incomes

Note:
optional `search`, `income_date`, `min_amount`, and `max_amount`.

**GET**
```http
http://localhost:8080/incomes/list/1?search=bonus&income_date=2026-10-28
```

## 5.6. paginate incomes


**GET**
```http
http://localhost:8080/incomes/paginate/1?page=1&order=id&direction=asc
```

----------------------------
# 6. expenses
## 6.1. create expense

Note:
stored amount is `item_price` times `qty`. The wallet balance must cover it. `UpdateItemPrice` true also sets the item price.

**POST**
```http
http://localhost:8080/expenses
```
[*raw*]

```json
{
    "expense_date": "2026-10-28",
    "note": "morning coffee",
    "WalletId": 1,
    "item_id": 1,
    "qty": 1,
    "item_price": "500",
    "UpdateItemPrice": false
}
```


## 6.2. show expense


**GET**
```http
http://localhost:8080/expenses/1
```

## 6.3. update expense

Note:
the wallet balance must still cover any increase.

**PUT**
```http
http://localhost:8080/expenses/1
```
[*raw*]

```json
{
    "expense_date": "2026-10-29",
    "note": "evening coffee",
    "WalletId": 1,
    "item_id": 1,
    "qty": 2,
    "item_price": "500",
    "UpdateItemPrice": false
}
```


## 6.4. delete expense

Note:
adds the amount back to the wallet.

**DELETE**
```http
http://localhost:8080/expenses/1
```

## 6.5. list expenses

Note:
optional `search` (note), `expense_date`, `min_amount`, and `max_amount`.

**GET**
```http
http://localhost:8080/expenses/list/1?search=evening&expense_date=2026-10-29
```

## 6.6. paginate expenses


**GET**
```http
http://localhost:8080/expenses/paginate/1?page=1&order=id&direction=desc
```

----------------------------
# 7. cost-perspectives
## 7.1. create cost perspective

Note:
`type` is `percentage`, `ratio`, `rate_hour`, `rate_day`, or `rate_month`. Each category is `expense`, `income`, or `daily_expense`, stored as its own row.

**POST**
```http
http://localhost:8080/cost-perspectives
```
[*raw*]

```json
{
    "name": "markup",
    "value": "15",
    "type": "percentage",
    "applications": [
        {
            "money_insight_category": ["expense", "daily_expense"],
            "starting_value": "1000"
        }
    ]
}
```


## 7.2. show cost perspective


**GET**
```http
http://localhost:8080/cost-perspectives/1
```

## 7.3. update cost perspective

Note:
replaces the previous applications.

**PUT**
```http
http://localhost:8080/cost-perspectives/1
```
[*raw*]

```json
{
    "name": "fees",
    "value": "20",
    "type": "ratio",
    "applications": [
        {
            "money_insight_category": ["income"],
            "starting_value": "2000"
        }
    ]
}
```


## 7.4. delete cost perspective


**DELETE**
```http
http://localhost:8080/cost-perspectives/1
```

## 7.5. list cost perspectives

Note:
optional `name`, `type`, and `is_active` (`1` is true).

**GET**
```http
http://localhost:8080/cost-perspectives?name=fees&type=ratio
```

## 7.6. paginate cost perspectives


**GET**
```http
http://localhost:8080/cost-perspectives/paginate?name=fees&page=1&order=id&direction=asc
```

## 7.7. toggle cost perspective


**POST**
```http
http://localhost:8080/is-active/cost-perspectives/1?is_active=0
```

----------------------------
