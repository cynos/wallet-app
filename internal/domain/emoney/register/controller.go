package register

import (
	"github.com/gin-gonic/gin"

	"github.com/wallet-app/internal/domain/emoney/users"
	"github.com/wallet-app/internal/pkg/common/http/response"
)

type HTTPController struct {
	usersUseCase users.UseCase
}

func NewHTTPController(usersUseCase users.UseCase) *HTTPController {
	return &HTTPController{
		usersUseCase: usersUseCase,
	}
}

func (controller *HTTPController) Register(c *gin.Context) {
	if c.Request.Method != "POST" {
		response.Response(c, 400, false, "Unsupported http method", nil)
		return
	}

	username, password, ok := c.Request.BasicAuth()
	if !ok {
		response.Response(c, 400, false, "Invalid parameters, invalid usersname or password", nil)
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")

	if username == "" || password == "" || email == "" || name == "" {
		response.Response(c, 400, false, "Invalid parameters", nil)
		return
	}

	_, err := controller.usersUseCase.GenerateUser(c.Request.Context(), users.DTOUsers{
		Name:     name,
		Email:    email,
		Username: username,
		Password: password,
	})
	if err != nil {
		response.Response(c, 400, false, "Internal service error", gin.H{"error": err.Error()})
		return
	}

	response.Response(c, 200, true, "Success create users", nil)
}
