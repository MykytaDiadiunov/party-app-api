package resources

import "go-rest-api/internal/domain"

type SessionDto struct {
	Token string  `json:"token"`
	User  UserDto `json:"user"`
}

func (s SessionDto) DomainToDto(token string, user domain.User) SessionDto {
	u := UserDto{}
	return SessionDto{
		Token: token,
		User:  u.DomainToDto(user),
	}
}
