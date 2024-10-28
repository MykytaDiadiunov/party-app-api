package app

import (
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/database/repositories"
)

type likeService struct {
	likeRepo    repositories.LikeRepository
	userService UserService
}

type LikeService interface {
	Save(domainLike domain.Like) error
	FindByLikedId(likedId uint64) ([]domain.User, error)
	FindByLikerId(likerId uint64) ([]domain.User, error)
	Delete(domainLike domain.Like) error
	Exists(domainLike domain.Like) error
}

func NewLikeService(likeRepo repositories.LikeRepository, userService UserService) LikeService {
	return likeService{
		likeRepo:    likeRepo,
		userService: userService,
	}
}

func (l likeService) Save(domainLike domain.Like) error {
	err := l.likeRepo.Save(domainLike)
	if err != nil {
		return err
	}

	return nil
}

func (l likeService) FindByLikedId(likedId uint64) ([]domain.User, error) {
	likes, err := l.likeRepo.FindByLikedId(likedId)
	if err != nil {
		return []domain.User{}, err
	}

	users := []domain.User{}
	for _, like := range likes {
		user, err := l.userService.FindById(like.LikerId)
		if err != nil {
			return []domain.User{}, nil
		}
		users = append(users, user)
	}

	return users, nil
}

func (l likeService) FindByLikerId(likerId uint64) ([]domain.User, error) {
	likes, err := l.likeRepo.FindByLikerId(likerId)
	if err != nil {
		return []domain.User{}, err
	}

	users := []domain.User{}
	for _, like := range likes {
		user, err := l.userService.FindById(like.LikedId)
		if err != nil {
			return []domain.User{}, nil
		}
		users = append(users, user)
	}

	return users, nil
}

func (l likeService) Delete(domainLike domain.Like) error {
	err := l.likeRepo.Delete(domainLike)
	if err != nil {
		return err
	}

	return nil
}

func (l likeService) Exists(domainLike domain.Like) error {
	err := l.likeRepo.Exists(domainLike)
	if err != nil {
		return err
	}

	return nil
}
