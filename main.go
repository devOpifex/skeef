package main

import (
	"log"
	"net/http"
	"os"

	"github.com/devOpifex/opiflex/app"
	"github.com/devOpifex/opiflex/config"
	"github.com/dghubble/oauth1"
	tw "github.com/dghubble/oauth1/twitter"
)

func main() {
	config, err := config.Read()

	if err != nil {
		log.Fatal(err)
	}

	var callbackURL = "http://localhost:" + config.Port + "/" + config.TwitterCallbackPath

	twitterConfig := &oauth1.Config{
		ConsumerKey:    config.TwitterConsumerKey,
		ConsumerSecret: config.TwitterConsumerSecret,
		CallbackURL:    callbackURL,
		Endpoint:       tw.AuthorizeEndpoint,
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &app.Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Config:   config,
		Oauth:    *twitterConfig,
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
