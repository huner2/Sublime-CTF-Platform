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
			"email varchar(320) UNIQUE," +
			"admin integer NOT NULL DEFAULT '0'," +
			"id SERIAL PRIMARY KEY" +
			");")
	return err
}

func (db *ctfDB) queryUser(uname string) *userT {
	result, err := db.Query("SELECT username, email, admin, id from USERS where username=?;", uname)
	if err != nil {
		log.Println("Unable to query user " + uname + " with error: " + err.Error())
		return nil
	}
	result.Next()
	if result.Err() != nil {
		log.Println("Unable to query user " + uname + " with error: " + result.Err().Error())
		return nil
	}
	var user userT
	if err := result.Scan(&user.username, &user.email, &user.admin, &user.id); err != nil {
		log.Println("Unable to query user " + uname + " with error: " + err.Error())
		return nil
	}
	return &user
}

func (db *ctfDB) userExists(uname string) bool {
	result, err := db.Query("SELECT COUNT(id) from USERS where username=$1;", uname)
	if err != nil {
		log.Println("Unable to count users with username " + uname + " with error: " + err.Error())
		return true
	}
	result.Next()
	if result.Err() != nil {
		log.Println("Unable to count users with username " + uname + " with error: " + result.Err().Error())
		return true
	}
	var count int
	if err := result.Scan(&count); err != nil {
		log.Println("Unable to count users with username " + uname + " with error: " + err.Error())
	}
	return count != 0
}

func (db *ctfDB) adminExists() bool {
	result, err := db.Query("SELECT COUNT(id) from USERS where admin=1;")
	if err != nil {
		log.Println("Unable to count admins with error: " + err.Error())
		return true
	}
	result.Next()
	if result.Err() != nil {
		log.Println("Unable to count admins with error: " + result.Err().Error())
		return true
	}
	var count int
	if err := result.Scan(&count); err != nil {
		log.Println("Unable to count admins with error: " + err.Error())
	}
	return count != 0
}

func (db *ctfDB) createUser(uname string, salt string, hash string, email string, admin int) error {
	_, err := db.Exec("INSERT INTO USERS (username, salt, hash, email, admin) VALUES ($1, $2, $3, $4, $5);", uname, salt, hash, email, admin)
	if err != nil {
		log.Println("Unable to create user: " + uname + " with error: " + err.Error())
	}
	return err
}
