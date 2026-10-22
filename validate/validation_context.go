package validate

import "github.com/linn221/expenses/app"

// AppendFunc adds validation rules for the action being checked.
type AppendFunc = func(...Rule)

type ValidationContext struct {
	createValidateFunc       []func(*app.Container, *app.Session, AppendFunc)
	updateValidateFunc       []func(*app.Container, *app.Session, AppendFunc)
	deleteValidateFunc       []func(*app.Container, *app.Session, AppendFunc)
	toggleActiveValidateFunc []func(*app.Container, *app.Session, AppendFunc)
}

func (vc *ValidationContext) OnAny(f func(*app.Session, AppendFunc)) {
	vc.createValidateFunc = append(vc.createValidateFunc, func(c *app.Container, s *app.Session, appendFunc AppendFunc) {
		f(s, appendFunc)
	})
	vc.updateValidateFunc = append(vc.updateValidateFunc, func(c *app.Container, s *app.Session, appendFunc AppendFunc) {
		f(s, appendFunc)
	})
	vc.deleteValidateFunc = append(vc.deleteValidateFunc, func(c *app.Container, s *app.Session, appendFunc AppendFunc) {
		f(s, appendFunc)
	})
}

func (vc *ValidationContext) OnMutate(f func(*app.Session, AppendFunc)) {
	vc.createValidateFunc = append(vc.createValidateFunc, func(c *app.Container, s *app.Session, appendFunc AppendFunc) {
		f(s, appendFunc)
	})
	vc.updateValidateFunc = append(vc.updateValidateFunc, func(c *app.Container, s *app.Session, appendFunc AppendFunc) {
		f(s, appendFunc)
	})
}

func (vc *ValidationContext) OnMutateWithContainer(f func(*app.Container, *app.Session, AppendFunc)) {
	vc.createValidateFunc = append(vc.createValidateFunc, func(c *app.Container, s *app.Session, appendFunc AppendFunc) {
		f(c, s, appendFunc)
	})
	vc.updateValidateFunc = append(vc.updateValidateFunc, func(c *app.Container, s *app.Session, appendFunc AppendFunc) {
		f(c, s, appendFunc)
	})
}

func (vc *ValidationContext) OnDelete(f func(*app.Session, AppendFunc)) {
	vc.deleteValidateFunc = append(vc.deleteValidateFunc, func(c *app.Container, s *app.Session, appendFunc AppendFunc) {
		f(s, appendFunc)
	})
}

func (vc *ValidationContext) OnDeleteWithContainer(f func(*app.Container, *app.Session, AppendFunc)) {
	vc.deleteValidateFunc = append(vc.deleteValidateFunc, func(c *app.Container, s *app.Session, appendFunc AppendFunc) {
		f(c, s, appendFunc)
	})
}

func (vc *ValidationContext) OnActiveStatusChange(f func(*app.Session, AppendFunc)) {
	vc.toggleActiveValidateFunc = append(vc.toggleActiveValidateFunc, func(c *app.Container, s *app.Session, af AppendFunc) {
		f(s, af)
	})
}

func generateRules(vc *ValidationContext, c *app.Container, s *app.Session, action Action) []Rule {
	var rules []Rule
	appendFunc := func(added ...Rule) {
		rules = append(rules, added...)
	}
	switch action {
	case ActionCreate:
		for _, fx := range vc.createValidateFunc {
			fx(c, s, appendFunc)
		}
	case ActionUpdate:
		for _, fx := range vc.updateValidateFunc {
			fx(c, s, appendFunc)
		}
	case ActionDelete:
		for _, fx := range vc.deleteValidateFunc {
			fx(c, s, appendFunc)
		}
	case ActionActiveStatusChange:
		for _, fx := range vc.toggleActiveValidateFunc {
			fx(c, s, appendFunc)
		}
	}
	return rules
}
