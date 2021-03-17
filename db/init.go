package db

import (
	"database/sql"
	"log"
	"os"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Con *sql.DB
}

// DBExists Whether the database exists
func DBExists() bool {
	_, err := os.Stat("skeef.db")
	return err == nil
}

// DBCreate Create Database
func DBCreate() error {

	_, err := os.Create("skeef.db")

	if err != nil {
		return err
	}

	return nil
}

// DBRemove Remove Database
func DBRemove() error {
	return os.Remove("skeef.db")
}

// DBConnect Connect to database
func DBConnect() *sql.DB {
	db, err := sql.Open("sqlite3", "skeef.db")

	if err != nil {
		log.Fatal("Could not connect to local database")
	}

	return db
}

// CreateUserTable Create user table
func (DB *Database) CreateUserTable() error {

	_, err := DB.Con.Query(`CREATE TABLE users (
		emails VARCHAR(50) NOT NULL PRIMARY KEY,
		password VARCHAR(60) NOT NULL,
		admin INTEGER
	);`)

	if err != nil {
		return err
	}

	return nil
}
