package app

import (
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/devOpifex/skeef-app/db"
	"github.com/golangcollege/sessions"
	"github.com/justinas/alice"
)

// Application Application object
type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Database db.Database
	Session  *sessions.Session
	Setup    bool
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	if !app.Setup {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	app.render(w, r, []string{"ui/html/home.page.tmpl"}, templateData{})

}

// Handlers Returns all routes
func (app *Application) Handlers() http.Handler {
	standardMiddleware := alice.New(secureHeaders)
	dynamicMiddleware := alice.New(app.Session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/setup", dynamicMiddleware.Then(http.HandlerFunc(app.setup)))
	mux.Get("/static/", app.static())

	return app.Session.Enable(standardMiddleware.Then(mux))
}
