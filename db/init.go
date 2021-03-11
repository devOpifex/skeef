package db

import (
	"database/sql"
	"os"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

func dbExists() bool {
	_, err := os.Stat("skeef.sqlite")
	return os.IsExist(err)
}

func dbCreate() error {

	if dbExists() {
		return nil
	}

	_, err := os.Create("skeef.db")

	if err != nil {
		return err
	}

	return nil
}

func dbCreateUserTable() error {

	db, err := sql.Open("sqlite3", "skeef.db")

	if err != nil {
		return err
	}

	defer db.Close()

	_, err = db.Query(`CREATE TABLE users (
		username VARCHAR(50) NOT NULL PRIMARY KEY,
		password VARCHAR(60) NOT NULL,
		admin INTEGER
	);`)

	if err != nil {
		return err
	}

	return nil
}
