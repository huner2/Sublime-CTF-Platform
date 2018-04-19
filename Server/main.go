package main

import (
    //"encoding/json"
    "net/http"
    "os"
    "time"

    "log"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

func main() {
	var entry = "../client/build/index.html"
	var port = "8080"

	r := mux.NewRouter()

	// api := r.PathPrefix("/api/v1/").Subrouter()
	// TODO api.HandleFunc("/some", someFunc).Methods("GET")

	// Serve static assets directly.
	r.PathPrefix("/static").Handler(http.FileServer(http.Dir("client/build/")))

	r.PathPrefix("/").HandlerFunc(IndexHandler(entry))

	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    "127.0.0.1:" + port,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}
