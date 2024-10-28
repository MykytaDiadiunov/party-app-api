package resources

import "go-rest-api/internal/domain"

type MemberDto struct {
	Id    uint64 `json:"id"`
	Name  string `json:"username"`
	Email string `json:"email"`
}

func (m MemberDto) DomainToDto(domainUser domain.User) MemberDto {
	return MemberDto{
		Id:    domainUser.Id,
		Name:  domainUser.Name,
		Email: domainUser.Email,
	}
}

func (m MemberDto) DomainToDtoCollection(domainUsers []domain.User) []MemberDto {
	result := make([]MemberDto, len(domainUsers))

	for i := range domainUsers {
		result[i] = m.DomainToDto(domainUsers[i])
	}

	return result
}

type MemberExistsDto struct {
	IsJoined bool `json:"isJoined"`
}

func (m MemberExistsDto) ResultToDto(result bool) MemberExistsDto {
	return MemberExistsDto{
		IsJoined: result,
	}
}
