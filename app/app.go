package app

import (
	"log"
	"net/http"

	"github.com/devOpifex/skeef-app/db"
)

// Application Application object
type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Database db.Database
	Setup    bool
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if !app.Setup {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	app.render(w, r, []string{"ui/html/home.page.tmpl"}, templateData{})

}

// Handlers Returns all routes
func (app *Application) Handlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/setup", app.setup)
	mux.Handle("/static/", app.static())
	return mux
}
