package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	db = open(os.Getenv("DB_DSN"))
}

func open(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if nil != err {
		panic(err)
	}
	return db
}
