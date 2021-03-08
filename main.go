package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devOpifex/skeef/app"
	"github.com/devOpifex/skeef/config"
)

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

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
		Addr:     ":" + config.Port,
		ErrorLog: errorLog,
		Handler:  app.Handlers(),
	}

	infoLog.Printf("Listening on http://localhost:%s", config.Port)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
