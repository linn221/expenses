package models

import (
	"context"
	"errors"
	"testing"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/validate"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCategoryNameIsUniquePerUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Category{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&Category{Name: "Food", HasUserId: HasUserId{UserId: 1}}).Error; err != nil {
		t.Fatal(err)
	}

	container := &app.Container{DB: db}
	session := &app.Session{UserId: 1, Ctx: context.Background()}

	duplicate := &Category{Name: "Food", HasUserId: HasUserId{UserId: 1}}
	err = validate.Validate(container, session, duplicate, validate.ActionCreate)
	if !errors.Is(err, validate.ErrClient) {
		t.Fatalf("create duplicate name: got %v", err)
	}

	session.UserId = 2
	otherUser := &Category{Name: "Food", HasUserId: HasUserId{UserId: 2}}
	if err := validate.Validate(container, session, otherUser, validate.ActionCreate); err != nil {
		t.Fatalf("same name for another user: %v", err)
	}

	session.UserId = 1
	session.ResId = 1
	ownName := &Category{Name: "Food", HasUserId: HasUserId{UserId: 1}}
	if err := validate.Validate(container, session, ownName, validate.ActionUpdate); err != nil {
		t.Fatalf("update keeping the same name: %v", err)
	}

	if err := db.Create(&Category{Name: "Rent", HasUserId: HasUserId{UserId: 1}}).Error; err != nil {
		t.Fatal(err)
	}
	session.ResId = 2
	clash := &Category{Name: "Food", HasUserId: HasUserId{UserId: 1}}
	err = validate.Validate(container, session, clash, validate.ActionUpdate)
	if !errors.Is(err, validate.ErrClient) {
		t.Fatalf("update onto an existing name: got %v", err)
	}
}
