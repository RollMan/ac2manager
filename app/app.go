package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RollMan/ac2manager/app/db"
	"github.com/RollMan/ac2manager/app/handlers"
	"github.com/RollMan/ac2manager/app/apiHandlers"

	"github.com/gorilla/mux"
)

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
  log.Print("Server started.")
	{
    dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/ac2?charset=utf8&parseTime=true", os.Getenv("AC2_DB_USERNAME"), os.Getenv("MYSQL_ROOT_PASSWORD"))
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

  // API for public
  r.HandleFunc("/api/races", apiHandlers.RacesHandler).Methods("GET")
  r.HandleFunc("/api/upcoming_race", apiHandlers.UpcomingRaceHandler).Methods("GET")
  r.HandleFunc("/api/past_races", apiHandlers.PastRacesHandler).Methods("GET")
  r.HandleFunc("/api/future_races", apiHandlers.FutureRacesHandler).Methods("GET")
  r.HandleFunc("/api/server_status", apiHandlers.ServerStatusHandler).Methods("GET")
  r.HandleFunc("/api/login", apiHandlers.LoginHandler).Methods("POST")

  // API requiring authentication
  r.HandleFunc("/api/add_race", handlers.AuthMiddleware(apiHandlers.AddRaceHandler)).Methods("POST")
  r.HandleFunc("/api/edit_race", handlers.AuthMiddleware(apiHandlers.EditRaceHandler)).Methods("POST")


	log.Fatal(http.ListenAndServe(":80", r))
}
