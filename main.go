package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devOpifex/skeef-app/app"
	"github.com/devOpifex/skeef-app/config"
)

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	config, err := config.Read()

	if err != nil {
		errorLog.Fatal(err)
	}

	app := &app.Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Config:   config,
	}

	srv := &http.Server{
		Addr:     config.Port,
		ErrorLog: errorLog,
		Handler:  app.Handlers(),
	}

	// check again every 30 minutes
	go func() {
		for range time.Tick(time.Minute * 30) {
			app.LicenseCheck()
		}
	}()

	infoLog.Printf("Listening on port%s", config.Port)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
