package app

import (
	"context"
	"errors"
)

type userIdKey struct{}
type sessionIdKey struct{}

func GetUserId(ctx context.Context) (int, error) {
	userId, ok := ctx.Value(userIdKey{}).(int)
	if !ok || userId <= 0 {
		return 0, errors.New("userId not found in context")
	}
	return userId, nil
}

func SetUserId(ctx context.Context, userId int) context.Context {
	return context.WithValue(ctx, userIdKey{}, userId)
}

func GetSessionId(ctx context.Context) (string, error) {
	sessionId, ok := ctx.Value(sessionIdKey{}).(string)
	if !ok || sessionId == "" {
		return "", errors.New("sessionId not found in context")
	}
	return sessionId, nil
}

func SetSessionId(ctx context.Context, sessionId string) context.Context {
	return context.WithValue(ctx, sessionIdKey{}, sessionId)
}
