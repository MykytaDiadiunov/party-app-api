package requests

import "go-rest-api/internal/domain"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r RegisterRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}, nil
}

func (r LoginRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Email:    r.Email,
		Password: r.Password,
	}, nil
}
