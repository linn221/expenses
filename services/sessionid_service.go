package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/linn221/expenses/app"
)

type sessionIdService struct {
	*app.Container
}

func (s *sessionIdService) NewSessionId(userId int) (string, error) {
	uuid := uuid.NewString()
	key := fmt.Sprintf("Session:%s", uuid)
	value := fmt.Sprint(userId)
	err := s.Cache.SetValue(key, value, time.Hour*24*7) // seven days
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (s *sessionIdService) RemoveSessionId(sessionId string) error {
	key := fmt.Sprintf("Session:%s", sessionId)
	return s.Cache.RemoveKey(key)
}

// func (s *SessionIdService) RemoveActiveSessions(userId int) error
