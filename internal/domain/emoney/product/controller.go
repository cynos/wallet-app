package product

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"

	"github.com/wallet-app/internal/infrastructure/cache"
	"github.com/wallet-app/internal/pkg/common/http/response"
)

type HTTPController struct {
	cache *cache.CacheRedis
}

func NewHTTPController(cache *cache.CacheRedis) *HTTPController {
	return &HTTPController{cache: cache}
}

func (controller *HTTPController) ProductList(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().Get("https://phoenix-imkas.ottodigital.id/interview/biller/v1/list")
	if err != nil || resp.StatusCode() != 200 {
		response.Response(c, 400, false, "Internal service error", gin.H{"error": err.Error()})
		return
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}

	response.Response(c, 200, true, "Record found", result)
}

func (controller *HTTPController) ProductDetail(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().Get("https://phoenix-imkas.ottodigital.id/interview/biller/v1/detail?billerId=" + c.Param("id"))
	if err != nil || resp.StatusCode() != 200 {
		response.Response(c, 400, false, "Internal service error", gin.H{"error": err.Error()})
		return
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		response.Response(c, 400, false, "Record not found", gin.H{"error": err.Error()})
		return
	}

	response.Response(c, 200, true, "Record found", result)
}
