package controllers

import (
	"errors"
	"go-rest-api/internal/app"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/requests"
	"go-rest-api/internal/infra/http/resources"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PartyController struct {
	partyService  app.PartyService
	memberService app.MemberService
}

func NewPartyController(partyServ app.PartyService, memberService app.MemberService) PartyController {
	return PartyController{
		partyService:  partyServ,
		memberService: memberService,
	}
}

func (p PartyController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creatorUser := r.Context().Value(UserKey).(domain.User)
		domainParty, err := requests.Bind(r, requests.CreatePartyRequest{}, domain.Party{})
		if err != nil {
			BadRequest(w, err)
			return
		}

		if domainParty.Price < 1 {
			BadRequest(w, errors.New("the price cannot be lower than 1"))
			return
		}

		if creatorUser.Points < domainParty.Price {
			BadRequest(w, errors.New("creator points less than price"))
			return
		}

		domainParty.CreatorId = creatorUser.Id

		domainParty, err = p.partyService.Save(domainParty)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		domainPartyMembers, err := p.memberService.FindByPartyId(domainParty.Id)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		memberDto := resources.MemberDto{}
		partyDto := resources.PartyWithMembersDto{}

		Success(w, partyDto.DomainPartyWithMembersToDto(domainParty, memberDto.DomainToDtoCollection(domainPartyMembers)))
	}
}

func (p PartyController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyId := chi.URLParam(r, "partyId")
		if partyId == "" {
			BadRequest(w, errors.New("invalid partyId"))
			return
		}
		numericPartyId, err := strconv.ParseUint(partyId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid partyId"))
			return
		}

		domainParty, err := p.partyService.FindById(numericPartyId)
		if err != nil {
			NotFound(w, err)
			return
		}

		domainPartyMembers, err := p.memberService.FindByPartyId(domainParty.Id)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		memberDto := resources.MemberDto{}
		partyDto := resources.PartyWithMembersDto{}
		Success(w, partyDto.DomainPartyWithMembersToDto(domainParty, memberDto.DomainToDtoCollection(domainPartyMembers)))
	}
}

func (p PartyController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyId := chi.URLParam(r, "partyId")
		if partyId == "" {
			BadRequest(w, errors.New("invalid partyId"))
			return
		}

		newPartyDomain, err := requests.Bind(r, requests.UpdatePartyRequest{}, domain.Party{})
		if err != nil {
			BadRequest(w, err)
			return
		}

		numericPartyId, err := strconv.ParseUint(partyId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid partyId"))
			return
		}

		_, err = p.partyService.FindById(numericPartyId)
		if err != nil {
			NotFound(w, err)
			return
		}

		newPartyDomain.Id = numericPartyId

		domainParty, err := p.partyService.Update(newPartyDomain)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		domainPartyMembers, err := p.memberService.FindByPartyId(domainParty.Id)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		memberDto := resources.MemberDto{}
		partyDto := resources.PartyWithMembersDto{}
		Success(w, partyDto.DomainPartyWithMembersToDto(domainParty, memberDto.DomainToDtoCollection(domainPartyMembers)))
	}
}

func (p PartyController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partyId := chi.URLParam(r, "partyId")
		if partyId == "" {
			BadRequest(w, errors.New("invalid partyId"))
			return
		}

		numericPartyId, err := strconv.ParseUint(partyId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid partyId"))
			return
		}

		err = p.partyService.Delete(numericPartyId)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (p PartyController) GetParties() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")
		if page == "" || limit == "" {
			BadRequest(w, errors.New("invalid page or limit"))
			return
		}
		numericPage, pErr := strconv.ParseInt(page, 10, 32)
		numericLimit, lErr := strconv.ParseInt(limit, 10, 32)

		if pErr != nil || lErr != nil {
			BadRequest(w, errors.New("invalid page or limit"))
			return
		}

		domainParties, err := p.partyService.GetParties(int32(numericPage), int32(numericLimit))
		if err != nil {
			NotFound(w, err)
			return
		}
		partyDto := resources.PartyDto{}
		Success(w, partyDto.DomainToDtoCollection(domainParties))
	}
}

func (p PartyController) FindByCreatorId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creatorId := chi.URLParam(r, "creatorId")
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")
		if creatorId == "" || page == "" || limit == "" {
			BadRequest(w, errors.New("invalid partyId"))
			return
		}
		numericCreatorId, err := strconv.ParseUint(creatorId, 10, 64)
		numericPage, pErr := strconv.ParseInt(page, 10, 32)
		numericLimit, lErr := strconv.ParseInt(limit, 10, 32)

		if err != nil || pErr != nil || lErr != nil {
			BadRequest(w, errors.New("invalid partyId"))
			return
		}

		domainParties, err := p.partyService.FindByCreatorId(numericCreatorId, int32(numericPage), int32(numericLimit))
		if err != nil {
			NotFound(w, err)
			return
		}
		partyDto := resources.PartyDto{}
		Success(w, partyDto.DomainToDtoCollection(domainParties))
	}
}
