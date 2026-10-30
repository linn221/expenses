package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/config"
)

type sessionIdService struct {
	*app.Container
}

func (s *sessionIdService) NewSessionId(userId int) (string, error) {
	uuid := uuid.NewString()
	key := fmt.Sprintf("Session:%s", uuid)
	value := fmt.Sprint(userId)
	err := s.Cache.SetValue(key, value, config.Cfg.SessionTTL())
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
