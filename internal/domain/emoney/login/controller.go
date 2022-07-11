package login

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/wallet-app/internal/domain/emoney/users"
	"github.com/wallet-app/internal/pkg/common/http/cookies"
	"github.com/wallet-app/internal/pkg/common/http/middleware"
	"github.com/wallet-app/internal/pkg/common/http/response"
	"github.com/wallet-app/internal/tools"
)

type HTTPController struct {
	usersUseCase users.UseCase
}

func NewHTTPController(usersUseCase users.UseCase) *HTTPController {
	return &HTTPController{
		usersUseCase: usersUseCase,
	}
}

func (controller *HTTPController) Login(c *gin.Context) {
	if c.Request.Method != "POST" {
		response.Response(c, 400, false, "Unsupported http method", nil)
		return
	}

	username, password, ok := c.Request.BasicAuth()
	if !ok {
		response.Response(c, 400, false, "Invalid parameters, invalid usersname or password", nil)
		return
	}

	user, err := controller.usersUseCase.AuthenticateUser(c, users.DTOUsers{
		Username: username,
		Password: password,
	})
	if err != nil {
		response.Response(c, 400, false, "Invalid username or password", nil)
		return
	}

	var JWT_SIG_KEY = os.Getenv("JWT_SIG_KEY")
	var JWT_TOKEN_EXP = os.Getenv("JWT_TOKEN_EXP")

	// return saved token
	if _, valid := middleware.IsTokenExpired(user.Token, JWT_SIG_KEY); valid {
		err = cookies.App.SecureCookie["auth"].SetValue(c.Writer, []byte(user.Token))
		if err != nil {
			response.Response(c, 400, false, "Internal service error", gin.H{"error": err.Error()})
			return
		}
		response.Response(c, 200, true, "Login Success", nil)
		return
	}

	var expires *jwt.NumericDate
	if rememberMe, _ := strconv.ParseBool(c.PostForm("rememberMe")); !rememberMe {
		// default expiration jwt = 1 week
		expires = jwt.NewNumericDate(time.Now().Add(time.Duration(tools.StringsToInt(JWT_TOKEN_EXP)) * time.Minute))
	} else {
		// set one month expired for jwt
		expires = jwt.NewNumericDate(time.Now().Add(time.Duration(43800) * time.Minute))
	}

	claims := struct {
		jwt.RegisteredClaims
		UserID   string `json:"UserID"`
		Username string `json:"Username"`
		Email    string `json:"Email"`
	}{
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: expires},
		UserID:           fmt.Sprint(user.ID),
		Username:         user.Username,
		Email:            user.Email,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(JWT_SIG_KEY))
	if err != nil {
		response.Response(c, 400, false, "Internal service error", gin.H{"error": err.Error()})
		return
	}

	// update token
	controller.usersUseCase.Update(c, users.Users{
		LastLogin: time.Now(),
		Token:     signedToken,
	}, int(user.ID))

	// set cookie auth
	err = cookies.App.SecureCookie["auth"].SetValue(c.Writer, []byte(signedToken))
	if err != nil {
		response.Response(c, 400, false, "Internal service error", gin.H{"error": err.Error()})
		return
	}

	response.Response(c, 200, true, "Login Success", nil)
}

func (controller *HTTPController) Logout(c *gin.Context) {
	cookies.App.SecureCookie["auth"].Delete(c.Writer)
	response.Response(c, 200, true, "Logout Success", nil)
}
