package app

import (
	"fmt"
	"net/http"

	"github.com/devOpifex/skeef-app/db"
)

func (app *Application) setup(w http.ResponseWriter, r *http.Request) {

	err := db.DBCreate()

	if err != nil {
		http.Error(w, "Failed to create database", http.StatusInternalServerError)
		return
	}

	app.Database = db.Database{Con: db.DBConnect()}
	err = app.Database.CreateUserTable()

	if err != nil {
		http.Error(w, "Failed to create user table", http.StatusInternalServerError)
		return
	}

	app.Setup = true

	fmt.Fprintf(w, "Setup here")
}
