package users

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/wallet-app/internal/infrastructure/cache"
	"github.com/wallet-app/internal/pkg/common/http/response"
	"github.com/wallet-app/internal/tools"
)

type HTTPController struct {
	usersUseCase UseCase
	cacher       cache.Cacher
}

func NewHTTPController(usersUseCase UseCase, cacher cache.Cacher) *HTTPController {
	return &HTTPController{
		usersUseCase: usersUseCase,
		cacher:       cacher,
	}
}

func (controller *HTTPController) UserAccount(c *gin.Context) {
	userInfo := c.MustGet("userInfo").(jwt.MapClaims)
	userID := tools.StringsToInt(userInfo["UserID"].(string))
	result, err := controller.usersUseCase.GetUserAccount(c, userID)
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}
	response.Response(c, 200, true, "Record found", gin.H{"data": result})
}

func (controller *HTTPController) UserBalance(c *gin.Context) {
	var userID int
	if c.Request.Header.Get("UserID") != "" {
		userID = tools.StringsToInt(c.Request.Header.Get("UserID"))
	} else {
		userInfo := c.MustGet("userInfo").(jwt.MapClaims)
		userID = tools.StringsToInt(userInfo["UserID"].(string))
	}
	result, err := controller.usersUseCase.GetUserBalance(c, userID)
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}
	response.Response(c, 200, true, "Record found", gin.H{"data": result})
}

func (controller *HTTPController) UserBalanceUpdate(c *gin.Context) {
	var params = struct {
		UsersID int `json:"users_id"`
		Amount  int `json:"amount"`
	}{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.Response(c, 400, false, "Invalid parameters", nil)
		return
	}

	err = controller.usersUseCase.SetUserBalance(c, params.UsersID, params.Amount)
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}

	response.Response(c, 200, true, "Success update balance", nil)
}

func (controller *HTTPController) UserTopups(c *gin.Context) {
	userInfo := c.MustGet("userInfo").(jwt.MapClaims)
	userID := tools.StringsToInt(userInfo["UserID"].(string))
	result, err := controller.usersUseCase.GetUserTopups(c, userID)
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}
	response.Response(c, 200, true, "Record found", gin.H{"data": result})
}

func (controller *HTTPController) UserPayments(c *gin.Context) {
	userInfo := c.MustGet("userInfo").(jwt.MapClaims)
	userID := tools.StringsToInt(userInfo["UserID"].(string))
	result, err := controller.usersUseCase.GetUserPayments(c, userID)
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}
	response.Response(c, 200, true, "Record found", gin.H{"data": result})
}
