package controllers

import (
	"errors"
	"go-rest-api/internal/app"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/resources"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MemberController struct {
	memberService app.MemberService
	partyServiece app.PartyService
}

func NewMemberController(memberServ app.MemberService, partyService app.PartyService) MemberController {
	return MemberController{
		memberService: memberServ,
		partyServiece: partyService,
	}
}

func (m MemberController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domainUser := r.Context().Value(UserKey).(domain.User)
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

		domainParty, err := m.partyServiece.FindById(numericPartyId)
		if err != nil {
			BadRequest(w, err)
			return
		}
		if domainParty.CreatorId == domainUser.Id {
			BadRequest(w, errors.New("owner can't join to party"))
			return
		}

		domainMember := domain.Member{
			PartyId: numericPartyId,
			UserId:  domainUser.Id,
		}

		err = m.memberService.Exists(domainMember)
		if err == nil {
			NoContent(w, errors.New("user already joined"))
			return
		}

		err = m.memberService.Save(domainMember)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (m MemberController) Exists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domainUser := r.Context().Value(UserKey).(domain.User)
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

		domainMember := domain.Member{
			PartyId: numericPartyId,
			UserId:  domainUser.Id,
		}

		err = m.memberService.Exists(domainMember)
		isExestsDto := resources.MemberExistsDto{}
		if err != nil {
			Success(w, isExestsDto.ResultToDto(false))
		} else {
			Success(w, isExestsDto.ResultToDto(true))
		}
	}
}

func (m MemberController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domainUser := r.Context().Value(UserKey).(domain.User)
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

		domainMember := domain.Member{
			PartyId: numericPartyId,
			UserId:  domainUser.Id,
		}

		err = m.memberService.Exists(domainMember)
		if err != nil {
			NoContent(w, err)
			return
		}

		err = m.memberService.Delete(domainMember)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
