package services

import "github.com/linn221/expenses/app"

type defaultService struct {
	*app.Container
}

// workaround since controller has no direct access to app container
func (s *defaultService) AccessContainer(userFunc func(*app.Container)) {
	userFunc(s.Container)
}
