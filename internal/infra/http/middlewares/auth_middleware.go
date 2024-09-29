package middlewares

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go-rest-api/internal/app"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/controllers"
	"net/http"
)

func AuthMiddleware(ja *jwtauth.JWTAuth, sessionServ app.SessionService, userServ app.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := jwtauth.VerifyRequest(ja, r, jwtauth.TokenFromHeader)
			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				controllers.Unauthorized(w, err)
				return
			}

			claims := token.PrivateClaims()
			userId := uint64(claims["user_id"].(float64))
			userUuid, err := uuid.Parse(claims["uuid"].(string))
			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			sess := domain.Session{
				UserId: userId,
				UUID:   userUuid,
			}
			err = sessionServ.Check(sess)
			if err != nil {
				controllers.Unauthorized(w, err)
				return
			}

			user, err := userServ.FindById(sess.UserId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					err = errors.New("token is unauthorized")
				}
				controllers.Unauthorized(w, err)
				return
			}
			ctx = context.WithValue(ctx, controllers.UserKey, user)
			ctx = context.WithValue(ctx, controllers.SessionKey, sess)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
