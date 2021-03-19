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
	License  db.License
}

type Setup struct {
	Tables  bool
	Admin   bool
	License bool
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	if !app.Database.AdminExists() {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	if !app.Database.LicenseExists() {
		http.Redirect(w, r, "/setup/validate", http.StatusSeeOther)
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
	mux.Get("/setup", dynamicMiddleware.Then(http.HandlerFunc(app.setupPage)))
	mux.Post("/setup", dynamicMiddleware.Then(http.HandlerFunc(app.setupForm)))
	mux.Get("/setup/validate", dynamicMiddleware.Then(http.HandlerFunc(app.validatePage)))
	mux.Post("/setup/validate", dynamicMiddleware.Then(http.HandlerFunc(app.validateForm)))
	mux.Get("/admin/signin", dynamicMiddleware.Then(http.HandlerFunc(app.signinPage)))
	mux.Post("/admin/signin", dynamicMiddleware.Then(http.HandlerFunc(app.signinForm)))
	mux.Get("/admin", dynamicMiddleware.Then(http.HandlerFunc(app.adminPage)))

	mux.Get("/static/", app.static())

	return app.Session.Enable(standardMiddleware.Then(mux))
}
