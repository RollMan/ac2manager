package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RollMan554/ac2manager/app/db"
	"github.com/RollMan554/ac2manager/app/handlers"

	"github.com/gorilla/mux"
)

func main() {
  log.Print("Server started.")
	{
		dsn := fmt.Sprintf("%s:%s@/ac2?charset=utf8", os.Getenv("AC2_DB_USERNAME"), os.Getenv("AC2_DB_PASSWORD"))
		db.InitDB(dsn)
    log.Print("DB OK.")
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RootHandler)
	r.HandleFunc("/login", handlers.LoginGetHandler).Methods("GET")
	r.HandleFunc("/login", handlers.LoginPostHandler).Methods("POST")
	r.HandleFunc("/admin", handlers.AuthMiddleware(handlers.AdminHandler))

	log.Fatal(http.ListenAndServe(":8000", r))
}
