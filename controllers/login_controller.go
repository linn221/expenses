package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/linn221/expenses/models"
	"github.com/linn221/expenses/services"
	"github.com/linn221/expenses/utils"
)

type loginController struct {
	Services *services.Provider
}

func (c *loginController) Login() http.HandlerFunc {
	return WithInputParsing(func(w http.ResponseWriter, r *http.Request, input *models.LoginInfo) error {

		var sessionId string
		user, err := c.Services.LoginService.Login(r.Context(), input)
		if err != nil {
			return err
		}

		sessionId, err = c.Services.SessionIdService.NewSessionId(user.Id)
		if err != nil {
			return err
		}

		utils.SetCookies(w, "session_id", sessionId)
		w.Write([]byte("ok"))
		return nil
	})
}

func (c *loginController) Register() http.HandlerFunc {
	return WithInputParsing(func(w http.ResponseWriter, r *http.Request, input *models.NewUser) error {
		result, err := c.Services.LoginService.Register(r.Context(), input)
		if err != nil {
			return err
		}
		res, err := getResource(c.Services.LoginService, result)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(res)
	})
}
