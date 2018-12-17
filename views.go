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

	"golang.org/x/crypto/blake2b"

	"github.com/flosch/pongo2"
)

var frame = pongo2.Must(pongo2.FromFile("./templates/frame.html")) // Only frame can be pre-compiled from what I can tell
var userRe = regexp.MustCompile(`[^[:alnum:]]`)
var emailRe = regexp.MustCompile(`.+\@.+\..+`)
var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func defaultContext(page string, user *userT, config *configT) *pongo2.Context {
	return &pongo2.Context{
		"title": config.ctfPrefs.title,
		"page":  page,
		"user": func() string {
			if user != nil {
				return user.username
			}
			return ""
		}(),
		"admin": func() bool {
			if user != nil {
				return user.admin == 1
			}
			return false
		}(),
		"login": user != nil,
	}
}

func indexView(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // Explicitly set content-type
	ctx := defaultContext("index.html", user, config)
	if err := frame.ExecuteWriter(*ctx, w); err != nil {
		log.Println("Unable to render index.html")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func challengeView(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx := defaultContext("challenges.html", user, config)
	if err := frame.ExecuteWriter(*ctx, w); err != nil {
		log.Println("Unable to render challenges.html")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func adminDashboardView(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	if user.admin != 1 {
		log.Println("Non-admin attempt to access dashboard.html")
		http.Error(w, "Not an admin", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx := defaultContext("dashboard.html", user, config)
	if err := frame.ExecuteWriter(*ctx, w); err != nil {
		log.Println("Unable to render dashboard.html")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func adminPagesView(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	if user.admin != 1 {
		log.Println("Non-admin attempt to access pages.html")
		http.Error(w, "Not an admin", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx := defaultContext("pages.html", user, config)
	if err := frame.ExecuteWriter(*ctx, w); err != nil {
		log.Println("Unable to render pages.html")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginView(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx := defaultContext("login.html", user, config)
	if err := frame.ExecuteWriter(*ctx, w); err != nil {
		log.Println("Unable to render login.html")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginSubmit(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
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
		log.Println("Invalid username in login request")
		w.Write([]byte("{\"success\": false}"))
		return
	}

	pword, ok := data["pword"].(string)
	if !ok {
		log.Println("Invalid password in login request")
		w.Write([]byte("{\"success\": false}"))
		return
	}

	uname = strings.TrimSpace(uname)
	info := config.db.loginUser(uname)
	if info == nil {
		log.Println("No user with that username")
		w.Write([]byte("{\"success\": false}"))
		return
	}

	rawhash := sha256.Sum256([]byte(info.salt + pword))
	hash := hex.EncodeToString(rawhash[:])
	if info.hash != hash {
		log.Println("Incorrect password")
		w.Write([]byte("{\"success\": false}"))
		return
	}

	created := time.Now().Unix()
	tokey := uname + string(created) + randStringBytes(16)
	rawkey := blake2b.Sum256([]byte(tokey))
	key := hex.EncodeToString(rawkey[:])
	err := config.db.createSession(info.id, created, key)
	if err != nil {
		log.Println("Unable to create session")
		w.Write([]byte("{\"success\": false}"))
	}

	w.Write([]byte("{\"success\": true, \"key\": \"" + key + "\"}"))
}

func registerSubmit(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
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

	uname = strings.TrimSpace(uname)
	email = strings.TrimSpace(email)

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
	if len(email) >= 320 || !emailRe.MatchString(email) {
		log.Println("Invalid email length")
		w.Write([]byte("{\"success\": false, \"error\": \"elen\"}"))
		return
	}

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

	uid, err := config.db.createUser(uname, salt, hash, email, admin)
	if err != nil {
		log.Println("Couldn't create user")
		http.Error(w, "Couldn't create user", http.StatusInternalServerError)
		return
	}

	created := time.Now().Unix()
	tokey := uname + string(created) + randStringBytes(16)
	rawkey := blake2b.Sum256([]byte(tokey))
	key := hex.EncodeToString(rawkey[:])
	err = config.db.createSession(uid, created, key)

	if err != nil {
		log.Println("Unable to create session")
		http.Error(w, "Couldn't create session", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{\"success\": true, \"key\": \"" + key + "\"}"))
}
