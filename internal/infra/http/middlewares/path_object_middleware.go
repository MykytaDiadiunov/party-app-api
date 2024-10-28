package middlewares

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-rest-api/internal/infra/http/controllers"
	"log"
	"net/http"
	"strconv"
)

type FindableT[T any] interface {
	Find(uint64) (T, error)
}

func PathObjectMiddleware[domainType any](service FindableT[domainType]) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			objectKey := controllers.ResolveCtxKeyFromPathType(new(domainType))
			objectId, err := strconv.ParseUint(chi.URLParam(r, objectKey.UrlParam()), 10, 64)
			if err != nil || objectId <= 0 {
				err = fmt.Errorf("invalid %s parameter(only positive uint)", objectKey.UrlParam())
				controllers.BadRequest(w, err)
				return
			}
			object, err := service.Find(objectId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					err = fmt.Errorf("record not found")
					log.Printf("service.Find(PathObjectMiddleware): %s", err)
					controllers.NotFound(w, err)
					return
				}
				log.Printf("service.Find(PathObjectMiddleware): %s", err)
				controllers.InternalServerError(w, err)
				return
			}
			ctx := controllers.GetPathValueInCtx(r.Context(), object)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}
