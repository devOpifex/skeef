package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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
func (DB *Database) CreateTableUser() error {

	_, err := DB.Con.Query(`CREATE TABLE users (
		email VARCHAR(50) NOT NULL PRIMARY KEY,
		hashed_password CHAR(60) NOT NULL,
		admin INTEGER
	);`)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) CreateTableLicense() error {

	_, err := DB.Con.Query(`CREATE TABLE license (
		email VARCHAR(50) NOT NULL PRIMARY KEY,
		license VARCHAR(255) NOT NULL
	);`)

	if err != nil {
		return err
	}

	return nil
}

// InsertUser Inserts a new user in the database
func (DB *Database) InsertUser(email, password string, admin int) error {

	stmt, err := DB.Con.Prepare("INSERT INTO users SET (user, hashed_password, admin) VALUES (?,?,?);")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(email, password, admin)

	if err != nil {
		return err
	}

	return nil
}

// InsertLicense Insert the user license
func (DB *Database) InsertLicense(email, license string) error {

	stmt, err := DB.Con.Prepare("INSERT INTO license SET (user, license) VALUES (?,?);")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(email, license)

	if err != nil {
		return err
	}

	return nil
}

func (DB *Database) Authenticate(email, password string) (string, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT email, hashed_password FROM licenses WHERE email = ?"
	row := DB.Con.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid credentials")
		} else {
			return "", err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", errors.New("invalid credentials")
		} else {
			return "", err
		}
	}

	// Otherwise, the password is correct. Return the user ID.
	return email, nil
}
