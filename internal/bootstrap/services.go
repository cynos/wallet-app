package bootstrap

import (
	"log"
	"time"

	"github.com/chmike/securecookie"
	"github.com/gin-gonic/gin"

	"github.com/wallet-app/internal/domain/emoney/login"
	"github.com/wallet-app/internal/domain/emoney/register"
	"github.com/wallet-app/internal/domain/emoney/users"
	"github.com/wallet-app/internal/domain/topup"
	"github.com/wallet-app/internal/pkg/common/http/cookies"
	"github.com/wallet-app/internal/pkg/common/http/middleware"
)

func serviceEmoney() {
	// init cookies
	sc := map[string]securecookie.Params{
		"auth": {
			Path:     "/",
			MaxAge:   (int(time.Minute) * 60) / int(time.Second),
			HTTPOnly: true,
			Secure:   false,
			SameSite: securecookie.Lax,
		},
	}
	if err := cookies.SetupCookies(sc); err != nil {
		log.Fatal(err.Error())
	}

	// migrate tables
	db.AutoMigrate(
		users.Users{},
	)

	// init repo
	usersRepo := users.NewRepository(db)
	// init usecase
	usersUseCase := users.NewUseCase(usersRepo)
	// init controller
	usersController := users.NewHTTPController(usersUseCase, cacher)
	loginController := login.NewHTTPController(usersUseCase)
	registerController := register.NewHTTPController(usersUseCase)

	// routes
	router := gin.Default()
	router.POST("/register", registerController.Register)
	router.POST("/login", loginController.Login)
	router.GET("/logout", loginController.Logout)
	{
		rUsers := router.Group("/users")
		rUsers.GET("/account", middleware.JWTAuthorization, usersController.UserAccount)
		rUsers.GET("/topupHistory", middleware.JWTAuthorization, usersController.UserTopupHistory)
		rUsers.POST("/balanceUpdate", usersController.UserBalanceUpdate)
	}

	router.Run()
}

func serviceTopUp() {
	// migrate tables
	db.AutoMigrate(
		topup.Topup{},
	)

	// init repo
	topupRepo := topup.NewRepository(db)
	// init usecase
	topupUsecase := topup.NewUseCase(topupRepo)
	// init controller
	topupController := topup.NewHTTPController(topupUsecase)

	// routes
	router := gin.Default()
	{
		routerTopup := router.Group("/topup")
		routerTopup.GET("/:id", topupController.GetByID)
		routerTopup.POST("/add", topupController.Add)
	}

	router.Run()
}
