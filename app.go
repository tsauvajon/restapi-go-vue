package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App : application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize : connects to postgresql
func (app *App) Initialize(user, password, dbname, host, port, sslmode string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", user, password, dbname, host, port, sslmode)

	var err error
	app.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter()
}

// Run : runs the app
func (app *App) Run(addr string) {

}
