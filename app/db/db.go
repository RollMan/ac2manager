package db

import (
	"database/sql"
	"fmt"
	"github.com/RollMan/ac2manager/app/models"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

var Db *sql.DB
var DbMap *gorp.DbMap

func InitDB(dsn string) {
	var err error
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("DB Opening Error.")
		log.Fatal(err)
	}
	for {
		err = Db.Ping()
		if err != nil {
			fmt.Println("DB Ping Error. Retrying...")
			log.Println(err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}
	DbMap = &gorp.DbMap{Db: Db, Dialect: gorp.MySQLDialect{}}
	DbMap.AddTableWithName(models.Event{}, "events")
	DbMap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
}
