package utils

import "github.com/gin-gonic/gin"

func Response(c *gin.Context, httpCode int, data interface{}) {
	c.JSON(httpCode, HandResponse(httpCode, data))
}

func HandResponse(httpCode int, data interface{}) interface{} {
	return gin.H{
		"err_no": httpCode,
		"data":   data,
	}
}
