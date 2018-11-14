package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var connString = "postgres://ctfu:ctuf@localhost/ctf"

type ctfDB struct {
	*sql.DB
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
