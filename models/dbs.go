// Package dbs 数据库操作类
package models

import (

	// nogolint

	"database/sql"
	"time"

	"github.com/Shuixingchen/year-service/utils"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var (
	// DBMaps 代表 sql 连接的 map
	DBMaps      map[string]*sql.DB
	SQLBatchNum = 100
)

// InitMySQLDB 初始化所有连接
func InitMySQLDB(c map[string]utils.MySQLDSN) {
	DBMaps = make(map[string]*sql.DB)
	for k, dsn := range c {
		tempDB, err := sql.Open("mysql", dsn.DSN)
		if err != nil {
			log.Error(err)
			err := tempDB.Close()
			if err != nil {
				log.Error(err)
				panic(err.Error())
			}
			panic(err.Error())
		}
		tempDB.SetConnMaxLifetime(time.Second * 60)
		tempDB.SetMaxIdleConns(30)
		tempDB.SetMaxOpenConns(50)
		DBMaps[k] = tempDB
	}
}

// CheckDBConnExists 检查连接是否存在
func CheckDBConnExists(conn string) bool {
	_, ok := DBMaps[conn]
	return ok
}
