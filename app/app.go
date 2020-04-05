package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RollMan/ac2manager/app/db"
	"github.com/RollMan/ac2manager/app/handlers"

	"github.com/gorilla/mux"
)

func main() {
  log.Print("Server started.")
	{
		dsn := fmt.Sprintf("%s:%s@/ac2?charset=utf8&parseTime=true", os.Getenv("AC2_DB_USERNAME"), os.Getenv("MYSQL_ROOT_PASSWORD"))
		db.InitDB(dsn)
    log.Print("DB OK.")
	}

	r := mux.NewRouter()

  // Public
	r.HandleFunc("/", handlers.RootHandler)
  r.HandleFunc("/about", handlers.AboutHandler)
	r.HandleFunc("/login", handlers.LoginGetHandler).Methods("GET")
	r.HandleFunc("/login", handlers.LoginPostHandler).Methods("POST")
	r.NotFoundHandler = http.StripPrefix("/", http.FileServer(http.Dir("static/")))

  // Need authentication
	r.HandleFunc("/admin", handlers.AuthMiddleware(handlers.AdminHandler)).Methods("GET")
  r.HandleFunc("/add", handlers.AuthMiddleware(handlers.AddHandler))
  r.HandleFunc("/add_event", handlers.AuthMiddleware(handlers.AddEventHandler))
  r.HandleFunc("/edit_event", handlers.AuthMiddleware(handlers.EditEventHandler)).Methods("GET")


	log.Fatal(http.ListenAndServe(":8000", r))
}
