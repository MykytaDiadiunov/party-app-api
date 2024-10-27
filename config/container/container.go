package container

import (
	"go-rest-api/internal/app"
	"go-rest-api/internal/infra/database"
	"go-rest-api/internal/infra/database/repositories"
	"go-rest-api/internal/infra/filesystem"
	"go-rest-api/internal/infra/http/controllers"
	"go-rest-api/internal/infra/http/middlewares"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

type Container struct {
	Services
	Controllers
	Middleware
}

type Services struct {
	app.UserService
	app.SessionService
	app.PartyService
	app.MemberService
}

type Controllers struct {
	controllers.UserController
	controllers.SessionController
	controllers.PartyController
	controllers.MemberController
}

type Middleware struct {
	AuthMw func(http.Handler) http.Handler
}

func New() Container {
	tknAuth := jwtauth.New("HS256", []byte("1234567890"), nil)
	db := database.New()

	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	partyRepo := repositories.NewPartyRepository(db)
	memberRepo := repositories.NewMemberRepository(db)

	imageService := filesystem.NewImageStorageService("file_storage")
	userService := app.NewUserService(userRepo)
	sessionService := app.NewSessionService(sessionRepo, userService, tknAuth)
	partyService := app.NewPartyService(partyRepo, imageService, userService)
	memberService := app.NewMemberService(memberRepo, userService, partyService)

	userController := controllers.NewUserController(userService)
	sessionController := controllers.NewSessionController(sessionService, userService)
	memberController := controllers.NewMemberController(memberService)
	partyController := controllers.NewPartyController(partyService, memberService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, sessionService, userService)

	return Container{
		Services: Services{
			userService,
			sessionService,
			partyService,
			memberService,
		},
		Controllers: Controllers{
			userController,
			sessionController,
			partyController,
			memberController,
		},
		Middleware: Middleware{
			authMiddleware,
		},
	}
}
