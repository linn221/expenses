package services_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPaginateUsesFilteredQuery(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Category{}); err != nil {
		t.Fatal(err)
	}

	rows := []models.Category{
		{Name: "rent", HasUserId: models.HasUserId{UserId: 1}},
		{Name: "food", HasUserId: models.HasUserId{UserId: 1}},
		{Name: "fuel", HasUserId: models.HasUserId{UserId: 1}},
		{Name: "other", HasUserId: models.HasUserId{UserId: 2}},
	}
	if err := db.Create(&rows).Error; err != nil {
		t.Fatal(err)
	}

	session := &app.Session{UserId: 1, Ctx: context.Background()}
	service := services.New(&app.Container{DB: db}).CategoryService
	filtered := service.ApplyFilters(session, &dtos.CategoryFilter{})

	page, err := services.Paginate(filtered, services.Pagination[models.Category]{
		Page:       2,
		PerPage:    2,
		Column:     "name",
		Direction:  "asc",
		Collection: dtos.RespondCategoryCollection,
	})
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 3 {
		t.Fatalf("total = %d, want 3", page.Total)
	}
	if page.Page != 2 || page.PerPage != 2 {
		t.Fatalf("page meta = %d/%d, want 2/2", page.Page, page.PerPage)
	}
	if len(page.Results) != 1 || page.Results[0].Name != "rent" {
		t.Fatalf("page results = %+v, want [rent]", page.Results)
	}
	encoded, err := json.Marshal(page.Data)
	if err != nil {
		t.Fatal(err)
	}
	var collection struct {
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	if err := json.Unmarshal(encoded, &collection); err != nil {
		t.Fatal(err)
	}
	if len(collection.Results) != 1 || collection.Results[0].Name != "rent" {
		t.Fatalf("collection = %s, want rent", encoded)
	}

	listed, err := service.List(session, &dtos.CategoryFilter{Name: "foo"})
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 || listed[0].Name != "food" {
		t.Fatalf("list results = %+v, want [food]", listed)
	}
}
