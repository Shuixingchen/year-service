package handlers

import (
	"database/sql"

	"github.com/Shuixingchen/year-service/models"
	"github.com/gin-gonic/gin"
)

type Years struct {
	DB *sql.DB
}

func NewYearsHandler() *Years {
	db := models.NewDB()
	return &Years{DB: db}
}

func (y *Years) SaveRecord(c *gin.Context) {
	var record models.Record
	err := c.ShouldBind(&record)
	if err != nil {
		c.String(500, "Valide params")
	}

}
