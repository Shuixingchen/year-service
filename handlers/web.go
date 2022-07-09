package handlers

import (
	"github.com/Shuixingchen/year-service/handlers/services"

	"github.com/gin-gonic/gin"
)

type WebHandler struct {
}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}
func (h *WebHandler) Handle() {
	r := gin.Default()
	years := services.NewYearsHandler()
	uniswapV3 := services.NewUniswapV3Handler()

	versionRoute := r.Group("/v1")
	serviceRoute := versionRoute.Group("/years")
	serviceRoute.POST("/record", years.SaveRecord)

	v3 := r.Group("/v3")
	uniswapV3Route := v3.Group("/uniswap")
	uniswapV3Route.POST("/swap", uniswapV3.Swap)
	r.Run(":8080")
}
