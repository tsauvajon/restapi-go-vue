package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App : application
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}

	if err := p.getProduct(app.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count < 1 || count > 10 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	products, err := getProducts(app.DB, start, count)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// defer : will be executed when the scope ends
	defer r.Body.Close()

	if err := p.createProduct(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadGateway, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.deleteProduct(app.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	entry := "client/dist/index.html"

	// open and parse a template text file
	// TODO: Refactor and clean
	if tpl, err := template.New("index").ParseFiles(entry); err != nil {
		log.Fatal(err)
	} else {
		tpl.Lookup("index").ExecuteTemplate(w, "index.html", nil)
	}
}

func (app *App) initializeRoutes() {
	static := "./client/dist/static/"

	api := app.Router.PathPrefix("/api/").Subrouter()

	api.HandleFunc("/products", app.getProducts).Methods("GET")
	api.HandleFunc("/products", app.createProduct).Methods("POST")
	api.HandleFunc("/products/{id:[0-9]+}", app.getProduct).Methods("GET")
	api.HandleFunc("/products/{id:[0-9]+}", app.updateProduct).Methods("PUT")
	api.HandleFunc("/products/{id:[0-9]+}", app.deleteProduct).Methods("DELETE")

	//app.Router.PathPrefix("/dist/").Handler(http.FileServer(http.Dir(static)))
	app.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))

	app.Router.HandleFunc("/", app.indexHandler).Methods("GET")
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
	app.initializeRoutes()
}

// Run : runs the app
func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
