package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func Init() {
	var err error
	db, err = sql.Open("mysql", "root:password@/flos?charset=utf8mb4&parseTime=true")
	if err != nil {
		panic(err)
	}
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}
