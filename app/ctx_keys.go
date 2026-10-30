package app

import (
	"context"
	"errors"
)

type userIdKey struct{}

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
