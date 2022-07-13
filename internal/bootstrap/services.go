package bootstrap

import (
	"log"
	"time"

	"github.com/chmike/securecookie"
	"github.com/gin-gonic/gin"

	"github.com/wallet-app/internal/domain/emoney/login"
	"github.com/wallet-app/internal/domain/emoney/payment"
	"github.com/wallet-app/internal/domain/emoney/product"
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
		payment.Payment{},
	)

	// init repo
	usersRepo := users.NewRepository(db)
	paymentRepo := payment.NewRepository(db)
	// init usecase
	usersUseCase := users.NewUseCase(usersRepo)
	paymentUseCase := payment.NewUseCase(paymentRepo)
	// init controller
	usersController := users.NewHTTPController(usersUseCase, cacher)
	loginController := login.NewHTTPController(usersUseCase)
	registerController := register.NewHTTPController(usersUseCase)
	productController := product.NewHTTPController(cacher)
	paymentController := payment.NewHTTPController(paymentUseCase)

	// routes
	router := gin.Default()
	{
		router.POST("/register", registerController.Register)
		router.POST("/login", loginController.Login)
		router.GET("/logout", loginController.Logout)
	}
	{
		rUsers := router.Group("/users")
		rUsers.GET("/account", middleware.JWTAuthorization, usersController.UserAccount)
		rUsers.GET("/balance", middleware.JWTAuthorization, usersController.UserBalance)
		rUsers.PUT("/balance", usersController.UserBalanceUpdate)
		rUsers.GET("/topups", middleware.JWTAuthorization, usersController.UserTopups)
		rUsers.GET("/payments", middleware.JWTAuthorization, usersController.UserPayments)
	}
	{
		router.GET("/payment/:trx", middleware.JWTAuthorization, paymentController.PaymentGet)
		router.POST("/payment", middleware.JWTAuthorization, paymentController.PaymentCreate)
		router.POST("/payment/confirm/:trx", middleware.JWTAuthorization, paymentController.PaymentConfirm)
	}
	{
		router.GET("/product", productController.ProductList)
		router.GET("/product/detail/:id", productController.ProductDetail)
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
