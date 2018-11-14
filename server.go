package main

import (
	"log"
	"net/http"
	"time"

	"github.com/flosch/pongo2"

	"github.com/gorilla/mux"
)

var frame = pongo2.Must(pongo2.FromFile("./templates/frame.html")) // Only frame can be pre-compiled from what I can tell

type ctfHandler struct {
	*configT
	handler func(config *configT, w http.ResponseWriter, r *http.Request)
}

func (ctfh *ctfHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctfh.handler(ctfh.configT, w, r)
}

func indexView(config *configT, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // Explicitly set content-type
	ctx := pongo2.Context{
		"title": config.ctfPrefs.title,
		"page":  "index.html",
	}
	if err := frame.ExecuteWriter(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	config, cerr := loadConfig()
	if cerr != nil {
		log.Fatal("Could not load config.ini")
	}

	// Init db
	db, dberr := openConnection()
	if dberr != nil {
		log.Fatal("Could not connect to db")
	}
	if dberr = db.createUserTable(); dberr != nil {
		log.Fatal("Could not create user table")
	}

	config.db = db

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Handle("/", &ctfHandler{config, indexView})

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
