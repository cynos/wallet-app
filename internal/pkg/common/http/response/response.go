package response

import "github.com/gin-gonic/gin"

func Response(c *gin.Context, status int, result bool, message string, ext interface{}) {
	var response = gin.H{
		"result":  result,
		"message": message,
	}

	if ext != nil {
		switch data := ext.(type) {
		case map[string]interface{}:
			for key, val := range data {
				response[key] = val
			}
		case gin.H:
			for key, val := range data {
				response[key] = val
			}
		case string:
			response["ext"] = data
		}
	}

	c.JSON(status, response)
}
