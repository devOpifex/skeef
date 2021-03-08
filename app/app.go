package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/devOpifex/skeef/config"
	"github.com/dghubble/oauth1"
)

// Application Application object
type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Config   config.Config
	Oauth    oauth1.Config
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.html",
		"./ui/html/base.layout.html",
		"./ui/html/footer.partial.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *Application) profile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are logged in")
}

// Handlers Returns all routes
func (app *Application) Handlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.Handle("/static/", app.static())
	mux.HandleFunc("/profile", app.profile)
	return mux
}
