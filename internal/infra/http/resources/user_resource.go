package resources

import "go-rest-api/internal/domain"

type UserDto struct {
	Id    uint64 `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
}

func (u UserDto) ToDomainModel(user domain.User) UserDto {
	return UserDto{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}
