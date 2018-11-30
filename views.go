package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/flosch/pongo2"
)

var frame = pongo2.Must(pongo2.FromFile("./templates/frame.html")) // Only frame can be pre-compiled from what I can tell
var userRe = regexp.MustCompile(`[^[:alnum:]]`)
var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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

func loginView(config *configT, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx := pongo2.Context{
		"title": config.ctfPrefs.title,
		"page":  "login.html",
	}
	if err := frame.ExecuteWriter(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerSubmit(config *configT, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var data map[string]interface{}
	if jerr := decoder.Decode(&data); jerr != nil {
		log.Println("Unable to decode data")
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
		return
	}

	uname, ok := data["uname"].(string)
	if !ok {
		log.Println("Invalid username in register request")
		w.Write([]byte("{\"success\": false, \"error\": \"invu\"}"))
		return
	}

	pword, ok := data["pword"].(string)
	if !ok {
		log.Println("Invalid password in register request")
		w.Write([]byte("{\"success\": false, \"error\": \"invp\"}"))
		return
	}

	email, ok := data["email"].(string)
	if !ok {
		log.Println("Invalid email in register request")
		w.Write([]byte("{\"success\": false, \"error\": \"inve\"}"))
		return
	}

	if len(uname) < 4 || len(uname) > 20 || userRe.MatchString(uname) {
		log.Println("Invalid username length")
		w.Write([]byte("{\"success\": false, \"error\": \"ulen\"}"))
		return
	}
	if len(pword) < 8 || len(pword) > 256 {
		log.Println("Invalid password length")
		w.Write([]byte("{\"success\": false, \"error\": \"plen\"}"))
		return
	}
	if len(email) >= 320 {
		log.Println("Invalid email length")
		w.Write([]byte("{\"success\": false, \"error\": \"elen\"}"))
		return
	}

	uname = strings.TrimSpace(uname)
	email = strings.TrimSpace(email)

	if config.db.userExists(uname) {
		log.Println("User already exists with name: " + uname)
		w.Write([]byte("{\"success\": false, \"error\": \"utake\"}"))
		return
	}

	rand.Seed(time.Now().UnixNano())
	salt := randStringBytes(16)
	rawhash := sha256.Sum256([]byte(salt + pword))
	hash := hex.EncodeToString(rawhash[:])
	admin := 0
	if !config.db.adminExists() {
		admin = 1
	}

	err := config.db.createUser(uname, salt, hash, email, admin)
	if err != nil {
		log.Println("Couldn't create user")
		http.Error(w, "Couldn't create user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{\"success\": true}"))
}
