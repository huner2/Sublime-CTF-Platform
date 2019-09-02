package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
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

func pageView(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	vars := mux.Vars(r)
	page := vars["page"]
	in, trueName := contains(config.pages, page)
	if !in {
		log.Println("Invalid page")
		http.Error(w, "Invalid page", http.StatusNotFound)
		return
	}
	ctx := defaultContext("../pages/"+trueName+".html", user, config)
	if err := frame.ExecuteWriter(*ctx, w); err != nil {
		log.Println("Unable to render page " + page)
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

func adminChallengeView(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	if user.admin != 1 {
		log.Println("Non-admin attempt to access adminchallenges.html")
		http.Error(w, "Not an admin", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=uft-8")
	ctx := defaultContext("adminchallenges.html", user, config)
	cats := config.db.getCats()
	challs := config.db.getChalls()
	catchall := make(map[string][]challT)
	for _, chall := range challs {
		for _, cat := range cats {
			if cat.id == chall.category {
				catchall[cat.name] = append(catchall[cat.name], chall)
			}
		}
	}
	uctx := pongo2.Context{
		"cats":     cats,
		"catchall": catchall,
	}
	ctx.Update(uctx)
	if err := frame.ExecuteWriter(*ctx, w); err != nil {
		log.Println("Unable to render adminchallenges.html")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateChallenge(user *userT, config *configT, w http.ResponseWriter, r *http.Request) {
	if user.admin != 1 {
		log.Println("Non-admin attempt to access adminchallenges.html")
		http.Error(w, "Not an admin", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var data map[string]interface{}
	if jerr := decoder.Decode(&data); jerr != nil {
		log.Println("Unable to decode data")
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
		return
	}
	ctype, ok := data["type"].(string)
	if !ok {
		log.Println("Unable to get type")
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}
	op, ok := data["operation"].(string)
	if !ok {
		log.Println("Unable to get operation")
		http.Error(w, "Invalid operation", http.StatusBadRequest)
		return
	}
	name, ok := data["name"].(string)
	if !ok {
		log.Println("Unable to get name")
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}
	if ctype == "category" {
		if op == "create" {
			if config.db.catExists(name) {
				log.Println("Category already exists")
				http.Error(w, "Category already exists", http.StatusBadRequest)
				return
			}
			if err := config.db.createCat(name); err != nil {
				log.Println("Unable to create category")
				http.Error(w, "Unable to create category", http.StatusInternalServerError)
				return
			}
		} else if op == "update" {
			if !config.db.catExists(name) {
				log.Println("Category doesn't exist")
				http.Error(w, "Category doesn't exist", http.StatusBadRequest)
				return
			}
			// TODO: More update code
		} else if op == "delete" {
			if err := config.db.deleteCat(name); err != nil {
				http.Error(w, "Unable to delete category", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Invalid operation", http.StatusBadRequest)
			return
		}
	} else if ctype == "challenge" {
		if op == "update" {
			desc, ok := data["desc"].(string)
			if !ok || len(desc) > 512 {
				http.Error(w, "Invalid description", http.StatusBadRequest)
				return
			}
			flag, ok := data["flag"].(string)
			if !ok || len(flag) == 0 {
				http.Error(w, "Invalid flag", http.StatusBadRequest)
				return
			}
			pointstr, ok := data["points"].(string)
			points, err := strconv.Atoi(pointstr)
			if !ok || err != nil || points < 0 {
				http.Error(w, "Invalid points amount", http.StatusBadRequest)
				return
			}
			cat, ok := data["cat"].(string)
			if !ok {
				http.Error(w, "Invalid category", http.StatusBadRequest)
				return
			}
			id, ok := data["id"].(float64)
			if !ok {
				http.Error(w, "Invalid id", http.StatusBadRequest)
				return
			}
			config.db.updateChallenge(int(id), cat, name, desc, flag, points)
		} else if op == "delete" {
			// TODO: Delete code
		} else {
			http.Error(w, "Invalid operation", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Invalid type", http.StatusBadRequest)
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
