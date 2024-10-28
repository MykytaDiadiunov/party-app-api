package resources

import "go-rest-api/internal/domain"

type UserDto struct {
	Id     uint64 `json:"id"`
	Name   string `json:"username"`
	Email  string `json:"email"`
	Points int32  `json:"points"`
}

func (u UserDto) DomainToDto(user domain.User) UserDto {
	return UserDto{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Points: user.Points,
	}
}

type UsersDto struct {
	Users []MemberDto `json:"users"`
}

func (u UsersDto) DomainToDtoCollection(domainUsers []domain.User) UsersDto {
	memberDto := MemberDto{}
	return UsersDto{
		Users: memberDto.DomainToDtoCollection(domainUsers),
	}
}
