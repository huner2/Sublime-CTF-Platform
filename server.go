package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type viewHandler struct {
	*configT
	authRequired bool
	handler      func(user *userT, config *configT, w http.ResponseWriter, r *http.Request)
}

func (vh *viewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key, _ := r.Cookie("key")
	var user *userT
	if key != nil {
		user = vh.configT.db.getSession(key.Value)
		if user == nil {
			w.Header().Add("Set-Cookie", "key=invalid; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT")
		}
	}
	if vh.authRequired && user == nil {
		http.Redirect(w, r, "/login", 307)
	}
	vh.handler(user, vh.configT, w, r)
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
	if dberr = db.createSessionTable(); dberr != nil {
		log.Fatal("Could not create session table with error: " + dberr.Error())
	}
	if dberr = db.createChallengeTables(); dberr != nil {
		log.Fatal("Could not create challenge tables with error: " + dberr.Error())
	}
	config.db = db

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Handle("/challenges", &viewHandler{config, true, challengeView}).Methods("GET")
	r.Handle("/login", &viewHandler{config, false, loginView}).Methods("GET")
	r.Handle("/login", &viewHandler{config, false, loginSubmit}).Methods("POST")
	r.Handle("/register", &viewHandler{config, false, registerSubmit}).Methods("POST")
	r.Handle("/admin", &viewHandler{config, true, adminDashboardView}).Methods("GET")
	r.Handle("/", &viewHandler{config, false, indexView})

	srv := &http.Server{
		Handler:      r,
		Addr:         strings.TrimSpace(config.ip) + ":" + strconv.Itoa(config.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
