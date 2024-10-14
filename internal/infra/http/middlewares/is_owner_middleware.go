package middlewares

import (
	"errors"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/http/controllers"
	"net/http"
)

type Usable interface {
	GetUserId() uint64
}

func IsOwnerMiddleware[domainType Usable]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user := ctx.Value(controllers.UserKey).(domain.User)
			object := controllers.GetPathValueFromCtx[domainType](ctx)

			if object.GetUserId() != user.Id {
				err := errors.New("you have no access to this object")
				controllers.Forbidden(w, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
