package db

import (
	"database/sql"
	"fmt"
	"github.com/RollMan/ac2manager/app/models"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func InitDB(dsn string) (*sql.DB, *gorp.DbMap) {
	var db *sql.DB
	var dbMap *gorp.DbMap

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("DB Opening Error.")
		log.Fatal(err)
	}
	for {
		err = db.Ping()
		if err != nil {
			log.Println("DB Ping Error. Retrying...")
			log.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("DB OK.")
		break
	}
	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	dbMap.AddTableWithName(models.Event{}, "events").SetKeys(true, "id")
	return db, dbMap
}
