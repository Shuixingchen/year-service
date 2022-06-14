package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() *sql.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/years?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}
