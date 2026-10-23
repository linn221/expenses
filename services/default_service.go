package services

import "github.com/linn221/expenses/app"

type defaultService struct {
	container *app.Container
}

// workaround since controller has no direct access to app container
// must go through this in order to acceess db and others, not letting direct access to app container
// generally not recommended. dependencies should always be injected
func (s *defaultService) AccessContainer(userFunc func(*app.Container)) {
	userFunc(s.container)
}
