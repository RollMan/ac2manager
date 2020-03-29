package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var Db *sql.DB

func InitDB(dsn string) {
	var err error
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Print("DB Opening Error")
		log.Fatal(err)
	}
	err = Db.Ping()
	if err != nil {
		fmt.Print("DB Ping Error")
		log.Fatal(err)
	}
}
