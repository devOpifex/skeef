package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devOpifex/skeef-app/app"
	"github.com/devOpifex/skeef-app/db"
	"github.com/golangcollege/sessions"
)

var session *sessions.Session
var secret = []byte("u46IpCV8y5Vlur8YvODJEhgOY8m9JVE5")

func main() {

	reset := flag.Bool("reset", false, "Reset the application and redo the first time setup.")
	addr := flag.String("addr", ":8080", "Address on which to serve the skeef.")
	flag.Parse()

	if *reset {
		db.DBRemove()
		return
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	session = sessions.New(secret)
	session.Lifetime = 12 * time.Hour
	session.SameSite = http.SameSiteStrictMode

	if !db.DBExists() {
		err := db.DBCreate()

		if err != nil {
			errorLog.Fatal("Could not create database")
			return
		}
	}

	app := &app.Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Setup: app.Setup{
			Tables:  false,
			Admin:   false,
			License: false,
		},
		Session:  session,
		Database: db.Database{Con: db.DBConnect()},
	}

	app.Setup.Tables = app.TablesExists()
	app.Setup.Admin = app.AdminExists()
	app.Setup.License = app.LicenseExists()

	defer app.Database.Con.Close()

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.Handlers(),
	}

	// check again every 30 minutes
	go func() {
		for range time.Tick(time.Minute * 30) {
			app.LicenseCheck()
		}
	}()

	infoLog.Printf("Listening on port%s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
