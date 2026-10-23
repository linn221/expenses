package validate

import (
	"context"
	"fmt"

	"github.com/linn221/expenses/app"
)

type existsRule struct {
	skipped  bool
	count    int64
	message  string
	table    string
	idColumn string
	filter   FilterFunc
	id       any
}

func ExistsRule(table string, id any, message string, filter FilterFunc) *existsRule {
	return &existsRule{
		message:  message,
		table:    table,
		id:       id,
		idColumn: "id",
		filter:   filter,
	}
}

func (r *existsRule) Skipped() bool {
	return r.skipped
}

func (r *existsRule) When(b bool) *existsRule {
	r.skipped = !b
	return r
}

func (r *existsRule) IdColumn(column string) *existsRule {
	r.idColumn = column
	return r
}

func (r *existsRule) Validate(ctx context.Context, c *app.Container) error {
	dbctx := c.DB.WithContext(ctx).Table(r.table)
	if r.filter != nil {
		r.filter(dbctx)
	}

	if err := dbctx.Where(fmt.Sprintf("%s = ?", r.idColumn), r.id).Count(&r.count).Error; err != nil {
		return serverError(err)
	}

	if r.count <= 0 {
		return clientErrorf("%s", r.message)
	}

	return nil
}

type uniqueRule struct {
	skipped  bool
	count    int64
	table    string
	column   string
	value    any
	exceptId int
	filter   FilterFunc
}

// exceptId works if > 0
func UniqueRule(table string, column string, value any, exceptId int, filter FilterFunc) *uniqueRule {
	return &uniqueRule{
		table:    table,
		column:   column,
		value:    value,
		exceptId: exceptId,
		filter:   filter,
	}
}

func (r *uniqueRule) When(b bool) *uniqueRule {
	r.skipped = !b
	return r
}

func (r *uniqueRule) Skipped() bool {
	return r.skipped
}

func (r *uniqueRule) Validate(ctx context.Context, c *app.Container) error {
	dbctx := c.DB.WithContext(ctx).Table(r.table).Where(r.column+" = ?", r.value)
	if r.filter != nil {
		r.filter(dbctx)
	}

	if r.exceptId > 0 {
		dbctx = dbctx.Where("id != ?", r.exceptId)
	}
	err := dbctx.Count(&r.count).Error
	if err != nil {
		return serverError(err)
	}
	if r.count > 0 {
		return clientErrorf("duplicate %s field value: %s", r.column, r.value)
	}

	return nil
}

type noResultRule struct {
	message string
	count   int64
	skipped bool
	table   string
	filter  FilterFunc
}

func NoResultRule(table string, message string, filter FilterFunc) *noResultRule {
	return &noResultRule{
		table:   table,
		message: message,
		filter:  filter,
	}
}

func (r *noResultRule) When(b bool) *noResultRule {
	r.skipped = !b
	return r
}

func (r *noResultRule) Skipped() bool {
	return r.skipped
}

func (r *noResultRule) Validate(ctx context.Context, c *app.Container) error {
	if r.filter == nil {
		panic("must apply filter for NoResultRule")
	}
	dbctx := c.DB.WithContext(ctx).Table(r.table)
	r.filter(dbctx)
	err := dbctx.Count(&r.count).Error
	if err != nil {
		return serverError(err)
	}
	if r.count > 0 {
		return clientError(r.message)
	}
	return nil
}

type failure struct {
	message string
}

// workaround for cases where a simple or complex if statement is what you need for validation, the Rule interface fails
func Failure(message string) *failure {
	f := failure{
		message: message,
	}

	return &f
}

func (r *failure) Skipped() bool {
	return false
}

func (r *failure) Validate(ctx context.Context, c *app.Container) error {
	return clientError(r.message)
}
