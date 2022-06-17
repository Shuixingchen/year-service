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

	versionRoute := r.Group("/v1")
	serviceRoute := versionRoute.Group("/years")
	serviceRoute.POST("/record", years.SaveRecord)
	r.Run(":8080")
}
