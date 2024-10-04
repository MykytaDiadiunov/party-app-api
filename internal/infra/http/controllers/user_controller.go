package controllers

import (
	"go-rest-api/internal/app"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/requests"
	"go-rest-api/internal/infra/http/resources"
	"net/http"
)

type UserController struct {
	userService app.UserService
}

func NewUserController(userService app.UserService) UserController {
	return UserController{userService: userService}
}

func (c UserController) UpdateMyBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		amount, err := requests.Bind(r, requests.UpdateMyBalanceRequest{}, domain.UpdateUserBalanceAmount{})
		if err != nil {
			BadRequest(w, err)
			return
		}

		updatedUser, err := c.userService.UpdateUserBalance(user, amount.Amount)
		if err != nil {
			BadRequest(w, err)
			return
		}
		Success(w, resources.UserDto{}.DomainToDto(updatedUser))
	}
}

func (c UserController) FindMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		Success(w, resources.UserDto{}.DomainToDto(user))
	}
}

func (c UserController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.RegisterRequest{}, domain.User{})
		if err != nil {
			BadRequest(w, err)
			return
		}

		user, err = c.userService.Save(user)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Created(w, user)
	}
}

func (c UserController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		err := c.userService.Delete(user.Id)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		Ok(w)
	}
}
