package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/linn221/expenses/dtos"
	"github.com/linn221/expenses/services"
	"github.com/linn221/expenses/utils"
)

type loginController struct {
	Services *services.Provider
}

func (c *loginController) Login() http.HandlerFunc {
	return WithInputParsing(func(w http.ResponseWriter, r *http.Request, input *dtos.LoginInfo) error {

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
		return json.NewEncoder(w).Encode(dtos.LoginResponse{SessionId: sessionId})
	})
}

func (c *loginController) Register() http.HandlerFunc {
	return WithInputParsing(func(w http.ResponseWriter, r *http.Request, input *dtos.NewUser) error {
		result, err := c.Services.LoginService.Register(r.Context(), input)
		if err != nil {
			return err
		}
		res := getResource(c.Services.LoginService, result)
		return json.NewEncoder(w).Encode(res)
	})
}
