package models

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

type Game struct {
	PlayAID   uint16
	PlayBID   uint16
	StartTime uint64
}

type Team struct {
	ID    uint16
	Name  string
	Icon  string
	Group string
}

func GetAllGames(db *sql.DB) []*Game {
	result := make([]*Game, 0)
	query := "select `playA_id`,`playB_id`, UNIX_TIMESTAMP(start_time) from games"
	rows, err := db.Query(query)
	if err != nil {
		log.WithField("method", "GetAllGames").Panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var g Game
		if err := rows.Scan(&g.PlayAID, &g.PlayBID, &g.StartTime); err == nil {
			result = append(result, &g)
		}
	}
	if err := rows.Err(); err != nil {
		log.WithField("method", "GetAllGames").Error(err)
	}
	return result
}

func GetAllTeams(db *sql.DB) map[uint16]*Team {
	result := make(map[uint16]*Team, 0)
	query := "select `id`,`name`,`icon`,`group` from teams"
	rows, err := db.Query(query)
	if err != nil {
		log.WithField("method", "GetAllGames").Panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Icon, &t.Group); err == nil {
			result[t.ID] = &t
		}
	}
	if err := rows.Err(); err != nil {
		log.WithField("method", "GetAllTeams").Error(err)
	}
	return result
}
