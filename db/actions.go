package db

import (
	"database/sql"
	"log"
	"os"
)

type Database struct {
	Con *sql.DB
}

// DBExists Whether the database exists
func Exists() bool {
	_, err := os.Stat("skeef.db")
	return err == nil
}

// DBCreate Create Database
func Create() error {
	_, err := os.Create("skeef.db")

	if err != nil {
		return err
	}

	return nil
}

// DBRemove Remove Database
func Remove() error {
	return os.Remove("skeef.db")
}

// DBConnect Connect to database
func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "skeef.db")

	if err != nil {
		log.Fatal("Could not connect to local database")
	}

	return db
}
