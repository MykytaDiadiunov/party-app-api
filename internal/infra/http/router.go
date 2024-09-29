package http

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go-rest-api/config/container"
	"go-rest-api/internal/infra/http/controllers"
	"net/http"
)

func CreateRouter(con container.Container) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Route("/ping", func(apiRouter chi.Router) {
			apiRouter.Get("/", pingHandler())
			apiRouter.Handle("/*", notFoundJson())
		})
		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, con.SessionController, con.AuthMw)
				})
				apiRouter.Route("/user", func(apiRouter chi.Router) {
					apiRouter.Use(con.AuthMw)
					UserRouter(apiRouter, con.UserController)
				})
			})
		})
	})

	return router
}

func AuthRouter(r chi.Router, sc controllers.SessionController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/register",
			sc.Register(),
		)
		apiRouter.Post(
			"/login",
			sc.Login(),
		)
		apiRouter.With(amw).Delete(
			"/logout",
			sc.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc controllers.UserController) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/me",
			uc.FindMe(),
		)
	})
}

func notFoundJson() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.NotFound(w, errors.New("resource Not Found"))
	}
}

func pingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.Ok(w)
	}
}
