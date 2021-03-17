package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devOpifex/skeef-app/app"
	"github.com/devOpifex/skeef-app/db"
)

func main() {

	addr := flag.String("addr", ":8080", "Address on which to serve the skeef")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	app := &app.Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Setup:    db.DBExists(),
	}

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
