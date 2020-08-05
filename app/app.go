package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RollMan/ac2manager/app/apiHandlers"
	"github.com/RollMan/ac2manager/app/db"
	"github.com/RollMan/ac2manager/app/handlers"

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


  // Requiring authentication
  r.PathPrefix("/admin/").HandlerFunc(handlers.AuthMiddleware(handlers.StaticWithAuth))


	// API for public
	r.HandleFunc("/api/races", apiHandlers.RacesHandler).Methods("GET")
	r.HandleFunc("/api/upcoming_race", apiHandlers.UpcomingRaceHandler).Methods("GET")
	r.HandleFunc("/api/past_races", apiHandlers.PastRacesHandler).Methods("GET")
	r.HandleFunc("/api/future_races", apiHandlers.FutureRacesHandler).Methods("GET")
	r.HandleFunc("/api/server_status", apiHandlers.ServerStatusHandler).Methods("GET")
	r.HandleFunc("/api/login", apiHandlers.LoginHandler).Methods("POST")

	// API requiring authentication
	r.HandleFunc("/api/add_race", handlers.AuthMiddleware(apiHandlers.AddRaceHandler)).Methods("POST")
	r.HandleFunc("/api/remove_race", handlers.AuthMiddleware(apiHandlers.RemoveRaceHandler)).Methods("POST")

	// Public
  r.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))

	log.Fatal(http.ListenAndServe(":80", r))
}
