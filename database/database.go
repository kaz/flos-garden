package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func Init() {
	var err error
	db, err = sqlx.Open("mysql", "root:password@/flos?charset=utf8mb4&parseTime=true")
	if err != nil {
		panic(err)
	}
}

func QueryRow(query string, args ...interface{}) *sqlx.Row {
	return db.QueryRowx(query, args...)
}
func Query(query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.Queryx(query, args...)
}
func Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}
