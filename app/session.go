package app

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type sessionValueKey int

// using enums to restrict the possibitlty of session values keys
const (
	SessionValueUpdateItemPrice sessionValueKey = iota
)

type Session struct {
	UserId int
	Ctx    context.Context
	ResId  int
	Tx     *gorm.DB
	values map[sessionValueKey]any
}

func (s *Session) SelectForUpdate(result any, where string, args ...any) error {
	if s.Tx == nil {
		panic("start a transaction before calling select with lock")
	}
	return s.Tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(where, args...).First(result).Error
}

func (s *Session) PutValue(key sessionValueKey, value any) {
	if s.values == nil {
		s.values = map[sessionValueKey]any{}
	}

	s.values[key] = value
}

func (s *Session) Retrieve(key sessionValueKey) any {
	value, ok := s.values[key]
	if !ok {
		// it should be a programmer error to retrieve from the Session with the key that does not exist
		panic("retrieving value from the Session that does not exist")
	}
	return value
}

func RetrieveFromSession[T any](s *Session, key sessionValueKey) T {
	value, ok := s.values[key]
	if !ok {
		// it should be a programmer error to retrieve from the Session with the key that does not exist
		panic("retrieving value from the Session that does not exist")
	}
	result, ok := value.(T)
	if !ok {
		panic("retrieving value from the Session with invalid type")
	}
	return result
}

func (s *Session) StartTransaction(c *Container) {
	tx := c.DB.WithContext(s.Ctx).Begin()
	s.Tx = tx
}

func (s *Session) Rollback() error {
	return s.Tx.Rollback().Error
}

func (s *Session) Commit() error {
	return s.Tx.Commit().Error
}
