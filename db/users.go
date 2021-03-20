package db

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// InsertUser Inserts a new user in the database
func (DB *Database) InsertUser(email, password string, admin int) error {

	stmt, err := DB.Con.Prepare("INSERT INTO users (email, hashed_password, admin) VALUES (?,?,?);")

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

func (DB *Database) Authenticate(email, password string) (string, error) {
	var hashedPassword []byte
	stmt := "SELECT hashed_password FROM users WHERE email = ?"
	row := DB.Con.QueryRow(stmt, email)
	err := row.Scan(&hashedPassword)

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
