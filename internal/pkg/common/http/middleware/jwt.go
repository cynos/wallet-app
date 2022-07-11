package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/wallet-app/internal/pkg/common/http/cookies"
	"github.com/wallet-app/internal/tools"
)

func JWTAuthorization(c *gin.Context) {
	var JWT_SIG_KEY = os.Getenv("JWT_SIG_KEY")
	var JWT_LOGIN_EXP = os.Getenv("JWT_LOGIN_EXP")

	tokenString, err := cookies.App.SecureCookie["auth"].GetValue(nil, c.Request)
	if err != nil {
		cookies.App.SecureCookie["auth"].Delete(c.Writer)
		c.Abort()
		c.String(400, "token expired or invalid")
		return
	}

	// prolongation expires cookie
	cookies.ReplenishExpiredCookie(
		"auth",
		(int(time.Minute)*tools.StringsToInt(JWT_LOGIN_EXP))/int(time.Second),
	).SetValue(c.Writer, tokenString)

	token, valid := IsTokenExpired(string(tokenString), JWT_SIG_KEY)
	if !valid {
		cookies.App.SecureCookie["auth"].Delete(c.Writer)
		c.Abort()
		c.String(400, "token expired or invalid")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		cookies.App.SecureCookie["auth"].Delete(c.Writer)
		c.Abort()
		c.String(400, "token expired or invalid")
		return
	}

	c.Set("userInfo", claims)
}

func IsTokenExpired(token, sigkey string) (*jwt.Token, bool) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("signing method invalid")
		}
		return []byte(sigkey), nil
	})

	if err != nil || !t.Valid {
		return t, false
	}

	return t, true
}
