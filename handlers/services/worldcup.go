package services

import (
	"database/sql"

	"github.com/Shuixingchen/year-service/models"
	"github.com/Shuixingchen/year-service/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WorldCup struct {
	DB *sql.DB
}

type GameInfo struct {
	ID           int
	PlayAID      uint16
	PlayBID      uint16
	PlayAIcon    string
	PlayBIcon    string
	PlayA        string
	PlayB        string
	StartTime    uint64
	ContractAddr string
}

func NewWorldCupHandler() *WorldCup {
	db := models.DBMaps["worldcup:write"]
	if db == nil {
		log.WithField("mysql", "db").Fatal("worldcup:write not exist")
	}
	return &WorldCup{DB: db}
}

func (h *WorldCup) GetAllGames(c *gin.Context) {
	result := make([]*GameInfo, 0)
	teams := models.GetAllTeams(h.DB)
	games := models.GetAllGames(h.DB)
	for k, val := range games {
		var g GameInfo
		g.ID = k + 1
		g.PlayAID = val.PlayAID
		g.PlayBID = val.PlayBID
		g.StartTime = val.StartTime
		g.ContractAddr = val.ContractAddr
		if t, ok := teams[val.PlayAID]; ok {
			g.PlayA = t.Name
			g.PlayAIcon = t.Icon
		}
		if t, ok := teams[val.PlayBID]; ok {
			g.PlayB = t.Name
			g.PlayBIcon = t.Icon
		}
		result = append(result, &g)
	}
	utils.Response(c, 200, result)
}
