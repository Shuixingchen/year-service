package handlers

import (
	"database/sql"
	"time"

	"github.com/Shuixingchen/year-service/handlers/common"
	"github.com/Shuixingchen/year-service/models"
	"github.com/gin-gonic/gin"
)

type Status = int

const (
	DefaultStatus = iota // default contract type
	ValidStatus
	OnChainStatus
	InvalidStatus
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
	record.CTime = uint64(time.Now().Unix())
	if common.VerifySign(record.Message, record.Signature, record.Address) {
		record.Status = ValidStatus
		models.SaveRecord(&record, y.DB)
		c.String(200, "success")
	} else {
		record.Status = InvalidStatus
		models.SaveRecord(&record, y.DB)
		c.String(500, "Valide Signature")
	}
}
