package resources

import (
	"go-rest-api/internal/domain"
	"time"
)

type PartyDto struct {
	Id          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       int32     `json:"price"`
	StartDate   time.Time `json:"startDate"`
	CreatorId   uint64    `json:"creatorId"`
}

func (p PartyDto) DomainToDto(domainParty domain.Party) PartyDto {
	return PartyDto{
		Id:          domainParty.Id,
		Title:       domainParty.Title,
		Description: domainParty.Description,
		Image:       domainParty.Image,
		Price:       domainParty.Price,
		StartDate:   domainParty.StartDate,
		CreatorId:   domainParty.CreatorId,
	}
}

type PartyWithMembersDto struct {
	Id          uint64      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Price       int32       `json:"price"`
	StartDate   time.Time   `json:"startDate"`
	CreatorId   uint64      `json:"creatorId"`
	Members     []MemberDto `json:"members"`
}

func (p PartyWithMembersDto) DomainPartyWithMembersToDto(domainParty domain.Party, members []MemberDto) PartyWithMembersDto {
	return PartyWithMembersDto{
		Id:          domainParty.Id,
		Title:       domainParty.Title,
		Description: domainParty.Description,
		Image:       domainParty.Image,
		Price:       domainParty.Price,
		StartDate:   domainParty.StartDate,
		CreatorId:   domainParty.CreatorId,
		Members:     members,
	}
}

type PartiesDto struct {
	Parties     []PartyDto `json:"items"`
	Total       uint64     `json:"total"`
	CurrentPage int32      `json:"currentPage"`
	LastPage    int32      `json:"lastPage"`
}

func (p PartyDto) DomainToDtoCollection(domainParties domain.Parties) PartiesDto {
	result := make([]PartyDto, len(domainParties.Parties))

	for i := range domainParties.Parties {
		result[i] = p.DomainToDto(domainParties.Parties[i])
	}

	return PartiesDto{
		Parties:     result,
		Total:       domainParties.Total,
		CurrentPage: domainParties.CurrentPage,
		LastPage:    domainParties.LastPage,
	}
}
