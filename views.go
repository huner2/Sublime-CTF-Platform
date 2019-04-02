package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

var frame = pongo2.Must(pongo2.FromFile("./templates/frame.html")) // Only frame can be pre-compiled from what I can tell
var userRe = regexp.MustCompile(`[^[:alnum:]]`)
var emailRe = regexp.MustCompile(`.+\@.+\..+`)
var pageRe = regexp.MustCompile(`([^a-zA-Z\d_-])`)
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
		"pages": config.pages,
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
	ctx := defaultContext("../pages/index.html", user, config)
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
	files, err := ioutil.ReadDir("./pages")
	if err != nil {
		log.Println("Unable to list pages")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	flist := make([]string, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".html") {
			flist = append(flist, strings.TrimSuffix(f.Name(), ".html"))
		}
	}
	(*ctx)["pages"] = flist
	if err := frame.ExecuteWriter(*ctx, w); err != nil {
		log.Println("Unable to render pages.html")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getPageSource(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	if user.admin != 1 {
		log.Println("Non-admin attempt to access pages.html")
		http.Error(w, "Not an admin", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	vars := mux.Vars(r)
	page := vars["page"]
	if pageRe.MatchString(page) {
		log.Println("Invalid page name")
		http.Error(w, "Invalid page name", http.StatusBadRequest)
		return
	}
	file, err := os.Open("./pages/" + page + ".html")
	if err != nil {
		log.Println("Unable to open the given page")
		http.Error(w, "Unknown Page", http.StatusBadRequest)
		return
	}
	io.Copy(w, file)
}

func updatePageSource(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	if user.admin != 1 {
		log.Println("Non-admin attempt to access pages.html")
		http.Error(w, "Not an admin", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	page := vars["page"]
	if pageRe.MatchString(page) {
		log.Println("Invalid page name")
		http.Error(w, "Invalid page name", http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var data map[string]interface{}
	if jerr := decoder.Decode(&data); jerr != nil {
		log.Println("Unable to decode data")
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
		return
	}
	operation, ok := data["operation"].(string)
	if !ok {
		log.Println("Unable to get operation")
		http.Error(w, "Invalid operation", http.StatusBadRequest)
		return
	}
	if operation == "create" {
		if _, err := os.Stat("./pages/" + page + ".html"); err == nil {
			log.Println("Page already exists")
			http.Error(w, "Page already exists", http.StatusBadRequest)
			return
		}
		_, err := os.Create("./pages/" + page + ".html")
		if err != nil {
			log.Println("Unable to create file for new page")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		config.pages = append(config.pages, page)
	} else if operation == "delete" {
		if page == "index" {
			log.Println("Cannot delete index")
			http.Error(w, "Cannot delete index", http.StatusBadRequest)
			return
		}
		for i, tpage := range config.pages {
			if tpage == page {
				config.pages = append(config.pages[:i], config.pages[i+1:]...)
			}
		}
		err := os.Remove("./pages/" + page + ".html")
		if err != nil {
			log.Println("Unable to delete page")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if operation == "update" {
		contents, ok := data["contents"].(string)
		if !ok {
			log.Println("Unable to get contents")
			http.Error(w, "Invalid contents", http.StatusBadRequest)
			return
		}
		f, err := os.OpenFile("./pages/"+page+".html", os.O_WRONLY|os.O_TRUNC, 644)
		if err != nil {
			log.Println("Unable to open given page")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = f.WriteString(contents)
		if err != nil {
			log.Println("Unable to write to file")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("Invalid operation for updatePage")
		http.Error(w, "Invalid operation", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
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
		return
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
