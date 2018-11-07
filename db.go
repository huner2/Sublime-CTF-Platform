package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var connString = "postgres://ctf:ctf@localhost/ctf"
var db *sql.DB

func openConnection() error {
	var err error
	db, err = sql.Open("postgres", connString)

	return err
}

// openConnection() must be called before this method
func createUserTable() error {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS user (" +
			"username varchar(20) UNIQUE," +
			"salt varchar(16) NOT NULL," +
			"hash varchar(32) NOT NULL," +
			"email varchar(320) UNIQUE," +
			"admin integer NOT NULL DEFAULT '0'," +
			"PRIMARY KEY (id)" +
			");")
	return err
}
