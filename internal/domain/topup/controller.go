package topup

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wallet-app/internal/pkg/common/http/response"
	"github.com/wallet-app/internal/tools"
)

type HTTPController struct {
	topupUseCase UseCase
}

func NewHTTPController(topupUseCase UseCase) *HTTPController {
	return &HTTPController{
		topupUseCase: topupUseCase,
	}
}

func (controller *HTTPController) GetByID(c *gin.Context) {
	id := c.Param("id")

	// get from db
	topup, err := controller.topupUseCase.GetByID(c, tools.StringsToInt(id))
	if err != nil {
		response.Response(c, 400, false, "Cannot get record", nil)
		return
	}

	response.Response(c, http.StatusOK, true, "Record found", gin.H{"data": topup})
}

func (controller *HTTPController) Add(c *gin.Context) {
	var err error
	var topup Topup

	err = c.Bind(&topup)
	if err != nil {
		response.Response(c, 400, false, "Invalid parameters", gin.H{"error": err.Error()})
		return
	}

	// update balance users
	err = controller.topupUseCase.UpdateBalance(c, int(topup.Amount), int(topup.UsersID))
	if err != nil {
		response.Response(c, 400, false, "Internal service error", gin.H{"error": err.Error()})
		return
	}

	_, err = controller.topupUseCase.Add(c.Request.Context(), topup)
	if err != nil {
		response.Response(c, 400, false, "Internal service error", gin.H{"error": err.Error()})
		return
	}

	response.Response(c, http.StatusOK, true, "Success add topup", nil)
}

func (controller *HTTPController) Update(c *gin.Context) {
	var err error
	var topup Topup

	err = c.Bind(&topup)
	if err != nil {
		response.Response(c, 400, false, "Invalid parameters", gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	_, err = controller.topupUseCase.Update(c.Request.Context(), topup, tools.StringsToInt(id))
	if err != nil {
		response.Response(c, 400, false, "Internal service error", gin.H{"error": err.Error()})
		return
	}

	response.Response(c, http.StatusOK, true, "Success update topup", nil)
}
