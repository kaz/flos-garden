package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func DB() *sqlx.DB {
	if db == nil {
		var err error
		db, err = sqlx.Open("mysql", "root:password@/flos?charset=utf8mb4&parseTime=true")
		if err != nil {
			panic(err)
		}
	}

	return db
}
