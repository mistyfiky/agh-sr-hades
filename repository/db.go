package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

func init() {
	defer func() {
		if result := recover(); nil != result {
			err, _ := result.(error)
			log.Fatal(err)
		}
	}()
	tmp, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if nil != err {
		panic(err)
	}
	db = tmp
}
