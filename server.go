package main

import (
	"log"
	"net/http"
	"time"

	"github.com/flosch/pongo2"

	"github.com/gorilla/mux"
)

var config configT
var frame = pongo2.Must(pongo2.FromFile("./templates/frame.html")) // Only frame can be pre-compiled from what I can tell

func indexView(w http.ResponseWriter, r *http.Request) {
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
	if cerr := loadConfig(); cerr != nil {
		log.Fatal("Could not load config.ini")
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.HandleFunc("/", indexView)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
