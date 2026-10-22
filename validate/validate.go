package validate

import (
	"github.com/linn221/expenses/app"
)

// State and Container will be injected by the service
// the Input struct will implement Input interface
func Validate(c *app.Container, session *app.Session, i Input, action Action) error {
	var vc ValidationContext
	i.ValidationRules(&vc)
	rules := generateRules(&vc, c, session, action)
	for _, rule := range rules {
		if rule.Skipped() {
			continue
		}
		err := rule.Validate(session.Ctx, c)
		if err != nil {
			return err
		}
	}

	return nil
}

// func ValidateInBatch(c *app.Container, i Input, state *Operation) error {
// 	var errMessage strings.Builder
// 	for _, rule := range i.Rules(state) {
// 		if rule.Skipped() {
// 			continue
// 		}
// 		if errStr := rule.Validate(c); errStr != "" {
// 			errMessage.Write([]byte(errStr + "\n"))
// 		}
// 	}

// 	if s := errMessage.String(); s != "" {
// 		return errors.New(s)
// 	}
// 	return nil
// }
