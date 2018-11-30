package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type viewHandler struct {
	*configT
	handler func(config *configT, w http.ResponseWriter, r *http.Request)
}

func (vh *viewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vh.handler(vh.configT, w, r)
}

func main() {
	config, cerr := loadConfig()
	if cerr != nil {
		log.Fatal("Could not load config.ini with error: " + cerr.Error())
	}

	// Init db
	db, dberr := openConnection()
	if dberr != nil {
		log.Fatal("Could not connect to db with error: " + dberr.Error())
	}
	if dberr = db.createUserTable(); dberr != nil {
		log.Fatal("Could not create user table with error: " + dberr.Error())
	}

	config.db = db

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Handle("/login", &viewHandler{config, loginView})
	r.Handle("/register", &viewHandler{config, registerSubmit}).Methods("POST")
	r.Handle("/", &viewHandler{config, indexView})

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
