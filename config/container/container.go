package container

import (
	"github.com/go-chi/jwtauth/v5"
	"go-rest-api/internal/app"
	"go-rest-api/internal/infra/database"
	"go-rest-api/internal/infra/database/repositories"
	"go-rest-api/internal/infra/http/controllers"
	"go-rest-api/internal/infra/http/middlewares"
	"net/http"
)

type Container struct {
	Services
	Controllers
	Middleware
}

type Services struct {
	app.UserService
	app.SessionService
}

type Controllers struct {
	controllers.UserController
	controllers.SessionController
}

type Middleware struct {
	AuthMw func(http.Handler) http.Handler
}

func New() Container {
	tknAuth := jwtauth.New("HS256", []byte("1234567890"), nil)
	db := database.New()

	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)

	userService := app.NewUserService(userRepo)
	sessionService := app.NewSessionService(sessionRepo, userService, tknAuth)

	userController := controllers.NewUserController(userService)
	sessionController := controllers.NewSessionController(sessionService, userService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, sessionService, userService)

	return Container{
		Services: Services{
			userService,
			sessionService,
		},
		Controllers: Controllers{
			userController,
			sessionController,
		},
		Middleware: Middleware{
			authMiddleware,
		},
	}
}
