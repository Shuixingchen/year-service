package models

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Record struct {
	Address     string `json:"address"`
	PayTxHash   string `json:"paytxhash"`
	Message     string `json:"message"`
	Signature   string `json:"signature"`
	Status      int
	BlockNumber uint64
	CTime       uint64
	UpdateTime  uint64
}

func SaveRecord(r *Record, db *sql.DB) {
	insertPrefix := "INSERT INTO records (address,pay_tx_hash,message,signature,status,block_number,c_time,update_time) VALUES "
	insertValue := fmt.Sprintf("('%s', '%s', '%s','%s','%d','%d','%d','%d')",
		r.Address,
		r.PayTxHash,
		r.Message,
		r.Signature,
		r.Status,
		r.BlockNumber,
		r.CTime,
		r.UpdateTime,
	)
	_, err := db.Exec(insertPrefix + insertValue)
	if err != nil {
		log.WithFields(log.Fields{"sql": insertPrefix + insertValue}).Error(err)
	}
}
