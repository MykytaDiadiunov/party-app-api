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

type LikeController struct {
	likeService app.LikeService
}

func NewLikeController(likeService app.LikeService) LikeController {
	return LikeController{
		likeService: likeService,
	}
}

func (l LikeController) SetLike() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		likerUser := r.Context().Value(UserKey).(domain.User)
		likedUserId := chi.URLParam(r, "likedId")

		if likedUserId == "" {
			BadRequest(w, errors.New("invalid likedId"))
			return
		}

		numericLikedUserId, err := strconv.ParseUint(likedUserId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid likedId"))
			return
		}

		if likerUser.Id == numericLikedUserId {
			BadRequest(w, errors.New("you can't like yourself"))
			return
		}

		domainLike := domain.Like{
			LikedId: numericLikedUserId,
			LikerId: likerUser.Id,
		}

		err = l.likeService.Exists(domainLike)
		if err == nil {
			BadRequest(w, errors.New("user already liked"))
			return
		}

		err = l.likeService.Save(domainLike)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (l LikeController) DeleteLike() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		likerUser := r.Context().Value(UserKey).(domain.User)
		likedUserId := chi.URLParam(r, "likedId")

		if likedUserId == "" {
			BadRequest(w, errors.New("invalid likedId"))
			return
		}

		numericLikedUserId, err := strconv.ParseUint(likedUserId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid likedId"))
			return
		}

		domainLike := domain.Like{
			LikedId: numericLikedUserId,
			LikerId: likerUser.Id,
		}

		err = l.likeService.Exists(domainLike)
		if err != nil {
			NoContent(w, err)
			return
		}

		err = l.likeService.Delete(domainLike)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (l LikeController) GetFavorites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		likerUser := r.Context().Value(UserKey).(domain.User)

		users, err := l.likeService.FindByLikerId(likerUser.Id)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		usersDto := resources.UsersDto{}
		Success(w, usersDto.DomainToDtoCollection(users))
	}
}

func (l LikeController) GetByLikedUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		likedUserId := chi.URLParam(r, "likedId")

		if likedUserId == "" {
			BadRequest(w, errors.New("invalid likedId"))
			return
		}

		numericLikedUserId, err := strconv.ParseUint(likedUserId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid likedId"))
			return
		}

		users, err := l.likeService.FindByLikedId(numericLikedUserId)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		usersDto := resources.UsersDto{}
		Success(w, usersDto.DomainToDtoCollection(users))
	}
}

func (l LikeController) LikeExists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		likerUser := r.Context().Value(UserKey).(domain.User)
		likedUserId := chi.URLParam(r, "likedId")

		if likedUserId == "" {
			BadRequest(w, errors.New("invalid likedId"))
			return
		}

		numericLikedUserId, err := strconv.ParseUint(likedUserId, 10, 64)
		if err != nil {
			BadRequest(w, errors.New("invalid likedId"))
			return
		}

		domainLike := domain.Like{
			LikedId: numericLikedUserId,
			LikerId: likerUser.Id,
		}

		err = l.likeService.Exists(domainLike)
		isExestsDto := resources.UserLikedDto{}
		if err != nil {
			Success(w, isExestsDto.ResultToDto(false))
		} else {
			Success(w, isExestsDto.ResultToDto(true))
		}

	}
}
