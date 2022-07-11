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
	user, err := controller.usersUseCase.GetUserAccount(c, tools.StringsToInt(userInfo["UserID"].(string)))
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}

	data := make(map[string]interface{})
	tools.ForEachStruct(user, func(key string, val interface{}) {
		if key != "Topups" {
			data[key] = val
		}
	})
	response.Response(c, 200, true, "Record found", gin.H{"data": data})
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

	user, err := controller.usersUseCase.GetUserAccount(c, params.UsersID)
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}

	user.Balance += int64(params.Amount)
	user, err = controller.usersUseCase.Update(c, user, int(user.ID))
	if err != nil {
		response.Response(c, 400, false, "Failed update balance", gin.H{"error": err.Error()})
		return
	}

	response.Response(c, 200, true, "Success update balance", nil)
}

func (controller *HTTPController) UserTopupHistory(c *gin.Context) {
	userInfo := c.MustGet("userInfo").(jwt.MapClaims)
	user, err := controller.usersUseCase.GetUserAccount(c, tools.StringsToInt(userInfo["UserID"].(string)))
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}
	response.Response(c, 200, true, "Record found", gin.H{"data": user.Topups})
}
