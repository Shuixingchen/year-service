package main

import (
	"github.com/Shuixingchen/year-service/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	years := handlers.NewYearsHandler()

	versionRoute := r.Group("/v1")
	versionRoute.POST("/record", years.SaveRecord)
	r.Run(":8080")
}
