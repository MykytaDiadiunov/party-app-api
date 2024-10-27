package app

import (
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/database/repositories"
)

type MemberService interface {
	Save(domainMember domain.Member) error
	Exists(domainMember domain.Member) error
	Delete(domainMember domain.Member) error
	FindByUserId(userId uint64) ([]domain.Party, error)
	FindByPartyId(partyId uint64) ([]domain.User, error)
}

type memberService struct {
	memberRepo   repositories.MemberRepository
	userService  UserService
	partyService PartyService
}

func NewMemberService(memberRepo repositories.MemberRepository, userService UserService, partyService PartyService) MemberService {
	return memberService{
		memberRepo:   memberRepo,
		userService:  userService,
		partyService: partyService,
	}
}

func (m memberService) Save(domainMember domain.Member) error {
	err := m.memberRepo.Save(domainMember)
	if err != nil {
		return err
	}

	return nil
}

func (m memberService) FindByUserId(userId uint64) ([]domain.Party, error) {
	members, err := m.memberRepo.FindByUserId(userId)
	if err != nil {
		return []domain.Party{}, err
	}

	parties := []domain.Party{}

	for _, member := range members {
		party, err := m.partyService.FindById(member.PartyId)
		if err != nil {
			return []domain.Party{}, err
		}
		parties = append(parties, party)
	}

	return parties, nil
}

func (m memberService) FindByPartyId(partyId uint64) ([]domain.User, error) {
	members, err := m.memberRepo.FindByPartyId(partyId)
	if err != nil {
		return []domain.User{}, err
	}

	users := []domain.User{}

	for _, member := range members {
		user, err := m.userService.FindById(member.UserId)
		if err != nil {
			return []domain.User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (m memberService) Delete(domainMember domain.Member) error {
	return m.memberRepo.Delete(domainMember)
}

func (m memberService) Exists(domainMember domain.Member) error {
	return m.memberRepo.Exists(domainMember)
}
