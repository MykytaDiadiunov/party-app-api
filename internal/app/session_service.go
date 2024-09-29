package app

import (
	"database/sql"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/database/repositories"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type SessionService interface {
	Register(user domain.User) (domain.User, string, error)
	Login(user domain.User) (domain.User, string, error)
	Logout(sess domain.Session) error
	Check(sess domain.Session) error
	GenerateToken(user domain.User) (string, error)
}

type sessionService struct {
	userServ    UserService
	sessionRepo repositories.SessionRepository
	tokenAuth   *jwtauth.JWTAuth
}

func NewSessionService(sr repositories.SessionRepository, us UserService, tokenAuth *jwtauth.JWTAuth) SessionService {
	return &sessionService{
		sessionRepo: sr,
		userServ:    us,
		tokenAuth:   tokenAuth,
	}
}

func (s sessionService) Register(user domain.User) (domain.User, string, error) {
	_, err := s.userServ.FindByEmail(user.Email)
	if err == nil {
		return domain.User{}, "", errors.New("invalid credentials")
	} else if !errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, "", err
	}

	user, err = s.userServ.Save(user)
	if err != nil {
		return domain.User{}, "", err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return domain.User{}, "", err
	}

	return user, token, nil
}

func (s sessionService) Login(user domain.User) (domain.User, string, error) {
	u, err := s.userServ.FindByEmail(user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("AuthService: failed to find user %s", err)
		}
		log.Printf("AuthService: login error %s", err)
		return domain.User{}, "", err
	}
	valid := s.checkPasswordHash(user.Password, u.Password)
	if !valid {
		return domain.User{}, "", errors.New("invalid credentials")
	}

	token, err := s.GenerateToken(u)
	if err != nil {
		return domain.User{}, "", err
	}

	return u, token, err
}

func (s sessionService) Logout(session domain.Session) error {
	return s.sessionRepo.Delete(session)
}

func (s sessionService) Check(session domain.Session) error {
	return s.sessionRepo.Exists(session)
}

func (s sessionService) GenerateToken(user domain.User) (string, error) {
	sess := domain.Session{UserId: user.Id, UUID: uuid.New()}
	err := s.sessionRepo.Save(sess)
	if err != nil {
		return "", err
	}

	claims := map[string]interface{}{
		"user_id": sess.UserId,
		"uuid":    sess.UUID,
	}
	jwtauth.SetExpiry(claims, time.Now().Add(72*time.Hour))

	_, tokenString, err := s.tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s sessionService) checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
