package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var connString = "postgres://ctfu:ctfu@localhost/ctf?sslmode=disable"

type ctfDB struct {
	*sql.DB
}

type userT struct {
	username string
	email    string
	admin    int
	id       int
}

type loginT struct {
	id   int
	hash string
	salt string
}

func openConnection() (*ctfDB, error) {
	var err error
	db, err := sql.Open("postgres", connString)

	ctfDB := &ctfDB{db}
	return ctfDB, err
}

// openConnection() must be called before this method
func (db *ctfDB) createUserTable() error {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS USERS (" +
			"username varchar(20) UNIQUE," +
			"salt varchar(16) NOT NULL," +
			"hash varchar(64) NOT NULL," +
			"email varchar(320) NOT NULL," +
			"admin integer NOT NULL DEFAULT '0'," +
			"id SERIAL PRIMARY KEY" +
			");")
	return err
}

// openConnection() must be called before this method
func (db *ctfDB) createSessionTable() error {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS SESSIONS (" +
			"id SERIAL," +
			"uid integer REFERENCES USERS(id)," +
			"created bigint NOT NULL," +
			"key varchar(64)," +
			"PRIMARY KEY (id,uid)" +
			");")
	return err
}

// openConnection() must be called before this method
// createSessionTable() must be called before this method
func (db *ctfDB) createSession(uid int, created int64, key string) error {
	_, err := db.Exec("DELETE FROM SESSIONS WHERE uid=$1;", uid)
	if err != nil {
		log.Println("Error deleting old sessions for uid " + string(uid) + " with error: " + err.Error())
	}
	_, err = db.Exec("INSERT INTO SESSIONS (uid, created, key) VALUES ($1, $2, $3);", uid, created, key)
	if err != nil {
		log.Println("Unable to create session for user: " + string(uid) + " with error: " + err.Error())
	}
	return err
}

// openConnection() must be called before this method
// createSessionTable() must be called before this method
// nil is returned if session is nil or invalid
func (db *ctfDB) getSession(session string) *userT {
	var user userT
	err := db.QueryRow("select USERS.username, USERS.email, USERS.admin, USERS.id from USERS INNER JOIN SESSIONS ON USERS.id = SESSIONS.uid WHERE SESSIONS.key=$1;", session).Scan(&user.username, &user.email, &user.admin, &user.id)
	if err != nil {
		log.Println("Unable to get session for key: " + session + " with error: " + err.Error())
		return nil
	}
	return &user
}

// openConnection() must be called before this method
// createUserTable() must be called before this method
func (db *ctfDB) queryUser(uname string) *userT {
	var user userT
	err := db.QueryRow("SELECT username, email, admin, id from USERS where LOWER(username)=LOWER($1);", uname).Scan(&user.username, &user.email, &user.admin, &user.id)
	if err != nil {
		log.Println("Unable to query user " + uname + " with error: " + err.Error())
		return nil
	}
	return &user
}

// openConnection() must be called before this method
// createUserTable() must be called before this method
func (db *ctfDB) loginUser(uname string) *loginT {
	var login loginT
	err := db.QueryRow("SELECT id, hash, salt from USERS where LOWER(username)=LOWER($1);", uname).Scan(&login.id, &login.hash, &login.salt)
	if err != nil {
		log.Println("Unable to query login information for user " + uname + " with error: " + err.Error())
		return nil
	}
	return &login
}

// openConnection() must be called before this method
// createUserTable() must be called before this method
func (db *ctfDB) userExists(uname string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(id) from USERS where LOWER(username)=LOWER($1);", uname).Scan(&count)
	if err != nil {
		log.Println("Unable to count users with username " + uname + " with error: " + err.Error())
		return true
	}
	return count != 0
}

// openConnection() must be called before this method
// createUserTable() must be called before this method
func (db *ctfDB) adminExists() bool {
	var count int
	err := db.QueryRow("SELECT COUNT(id) from USERS where admin=1;").Scan(&count)
	if err != nil {
		log.Println("Unable to count admins with error: " + err.Error())
		return true
	}
	return count != 0
}

// openConnection() must be called before this method
// createUserTable() must be called before this method
func (db *ctfDB) createUser(uname string, salt string, hash string, email string, admin int) (int, error) {
	var uid int
	err := db.QueryRow("INSERT INTO USERS (username, salt, hash, email, admin) VALUES ($1, $2, $3, $4, $5) RETURNING id;", uname, salt, hash, email, admin).Scan(&uid)
	if err != nil {
		log.Println("Unable to create user: " + uname + " with error: " + err.Error())
	}
	return uid, err
}
