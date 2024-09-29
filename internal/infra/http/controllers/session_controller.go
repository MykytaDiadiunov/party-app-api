package controllers

import (
	"errors"
	"go-rest-api/internal/app"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/requests"
	"go-rest-api/internal/infra/http/resources"
	"net/http"
)

type SessionController struct {
	sessionServ app.SessionService
	userServ    app.UserService
}

func NewSessionController(sessionServ app.SessionService, userServ app.UserService) SessionController {
	return SessionController{sessionServ, userServ}
}

func (c SessionController) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.RegisterRequest{}, domain.User{})
		if err != nil {
			BadRequest(w, errors.New("invalid request body"))
			return
		}

		user, token, err := c.sessionServ.Register(user)
		if err != nil {
			BadRequest(w, err)
			return
		}

		var sessDto resources.SessionDto
		Success(w, sessDto.DomainToDto(token, user))
	}
}

func (c SessionController) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domainUser, err := requests.Bind(r, requests.LoginRequest{}, domain.User{})
		if err != nil {
			BadRequest(w, errors.New("invalid request body"))
			return
		}
		user, token, err := c.sessionServ.Login(domainUser)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		var sessDto resources.SessionDto
		Success(w, sessDto.DomainToDto(token, user))
	}
}

func (c SessionController) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := r.Context().Value(SessionKey).(domain.Session)
		err := c.sessionServ.Logout(sess)
		if err != nil {
			InternalServerError(w, err)
		}
		Ok(w)
	}
}
