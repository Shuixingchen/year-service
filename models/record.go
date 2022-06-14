package models

import "database/sql"

type Record struct {
	Message   string `json:"message"`
	Address   string `json:"address"`
	TxHash    string `json:"txHash`
	Signature string `json:"signature"`
}

func SaveRecord(r *Record, db *sql.DB) {
	sql := "INSERT INTO record () "
}
