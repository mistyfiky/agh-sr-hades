package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	tmp, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if nil != err {
		panic(err)
	}
	db = tmp
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS `users` (" +
		"`username` VARCHAR(255) NOT NULL, " +
		"`password` VARCHAR(255) NOT NULL, " +
		"PRIMARY KEY (`username`))"); nil != err {
		panic(err)
	}
}
