package services

import (
	"context"
	"testing"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUpdateItemCategoryId(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Category{}, &models.Item{}); err != nil {
		t.Fatal(err)
	}

	food := models.Category{Name: "Food", HasUserId: models.HasUserId{UserId: 1}}
	rent := models.Category{Name: "Rent", HasUserId: models.HasUserId{UserId: 1}}
	if err := db.Create(&food).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&rent).Error; err != nil {
		t.Fatal(err)
	}

	item := models.Item{
		Name:       "latte",
		Code:       "latte",
		Price:      decimal.NewFromInt(4000),
		CategoryId: food.Id,
		HasUserId:  models.HasUserId{UserId: 1},
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	session := &app.Session{UserId: 1, ResId: item.Id, Ctx: context.Background()}
	service := &itemService{Container: &app.Container{DB: db}}

	existing, err := service.Get(session)
	if err != nil {
		t.Fatal(err)
	}

	input := *existing
	input.CategoryId = rent.Id
	updated, err := service.Update(session, existing, &input)
	if err != nil {
		t.Fatal(err)
	}

	var stored models.Item
	if err := db.First(&stored, item.Id).Error; err != nil {
		t.Fatal(err)
	}
	if stored.CategoryId != rent.Id {
		t.Fatalf("stored category_id = %d, want %d", stored.CategoryId, rent.Id)
	}
	if updated.Category.Name != rent.Name {
		t.Fatalf("response category = %q, want %q", updated.Category.Name, rent.Name)
	}
}

func TestUpdateItemCode(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Category{}, &models.Item{}); err != nil {
		t.Fatal(err)
	}

	food := models.Category{Name: "Food", HasUserId: models.HasUserId{UserId: 1}}
	if err := db.Create(&food).Error; err != nil {
		t.Fatal(err)
	}

	item := models.Item{
		Name:       "latte",
		Code:       "latte",
		Price:      decimal.NewFromInt(4000),
		CategoryId: food.Id,
		HasUserId:  models.HasUserId{UserId: 1},
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatal(err)
	}

	session := &app.Session{UserId: 1, ResId: item.Id, Ctx: context.Background()}
	service := &itemService{Container: &app.Container{DB: db}}
	existing, err := service.Get(session)
	if err != nil {
		t.Fatal(err)
	}

	input := dtos.NewItem{
		Name:       "latte",
		Code:       "mocha",
		Price:      decimal.NewFromInt(4000),
		CategoryId: food.Id,
	}.ToModel(session)
	updated, err := service.Update(session, existing, &input)
	if err != nil {
		t.Fatal(err)
	}

	var stored models.Item
	if err := db.First(&stored, item.Id).Error; err != nil {
		t.Fatal(err)
	}
	if stored.Code != "mocha" {
		t.Fatalf("stored code = %q, want %q", stored.Code, "mocha")
	}
	if updated.Code != "mocha" {
		t.Fatalf("response code = %q, want %q", updated.Code, "mocha")
	}
}
