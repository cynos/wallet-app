package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/wallet-app/internal/pkg/common/http/response"
	"github.com/wallet-app/internal/tools"
)

type HTTPController struct {
	paymentUseCase UseCase
}

func NewHTTPController(paymentUseCase UseCase) *HTTPController {
	return &HTTPController{
		paymentUseCase: paymentUseCase,
	}
}

func (controller *HTTPController) PaymentGet(c *gin.Context) {
	var (
		userInfo = c.MustGet("userInfo").(jwt.MapClaims)
		userID   = tools.StringsToInt(userInfo["UserID"].(string))
		trxid    = c.Param("trx")
	)

	res, err := controller.paymentUseCase.GetPayment(c, uint(userID), trxid)
	if err != nil {
		response.Response(c, 400, false, "Failed to get payment", gin.H{"error": err.Error()})
		return
	}
	response.Response(c, 200, true, "Success get payment", gin.H{"data": res})
}

func (controller *HTTPController) PaymentCreate(c *gin.Context) {
	userInfo := c.MustGet("userInfo").(jwt.MapClaims)
	userID := tools.StringsToInt(userInfo["UserID"].(string))

	var params = struct {
		ProductID string `json:"product_id"`
		Amount    int    `json:"amount"`
		Remark    string `json:"remark"`
	}{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.Response(c, 400, false, "Invalid parameters", nil)
		return
	}

	res, err := controller.paymentUseCase.CreatePayment(c, uint(userID), params.ProductID, params.Amount, params.Remark)
	if err != nil {
		response.Response(c, 400, false, "Failed to create payment", gin.H{"error": err.Error()})
		return
	}
	response.Response(c, 200, true, "Success create payment", gin.H{"data": res})
}

func (controller *HTTPController) PaymentConfirm(c *gin.Context) {
	var (
		userInfo = c.MustGet("userInfo").(jwt.MapClaims)
		userID   = tools.StringsToInt(userInfo["UserID"].(string))
		trxid    = c.Param("trx")
	)

	res, err := controller.paymentUseCase.ConfirmPayment(c, uint(userID), trxid)
	if err != nil {
		response.Response(c, 400, false, "Failed to confirm payment", gin.H{"error": err.Error()})
		return
	}
	response.Response(c, 200, true, "Success confirm payment", gin.H{"data": res})
}
