package app

import (
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/devOpifex/skeef-app/db"
	"github.com/devOpifex/skeef-app/graph"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/golangcollege/sessions"
	"github.com/justinas/alice"
)

// Application Application object
type Application struct {
	InfoLog         *log.Logger
	ErrorLog        *log.Logger
	Database        db.Database
	Session         *sessions.Session
	License         db.License
	Addr            string
	Count           int
	Stream          *twitter.Stream
	Valid           bool
	LicenseResponse LicenseResponse
	Pool            *Pool
	Quit            chan bool
	Start           chan bool
	Streaming       bool
	Graph           graph.Graph
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

	var tmplData templateData

	tmplData.Authenticated = app.isAuthenticated(r)

	tmplData.Streaming = app.Database.StreamOnGoing()

	app.render(w, r, []string{"ui/html/home.page.tmpl"}, tmplData)
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
	mux.Post("/admin", dynamicMiddleware.Then(http.HandlerFunc(app.adminForm)))
	mux.Get("/admin/signout", dynamicMiddleware.ThenFunc(app.signout))
	mux.Get("/ws", dynamicMiddleware.ThenFunc(app.socket))
	mux.Get("/admin/edit/:stream", dynamicMiddleware.Then(http.HandlerFunc(app.streamEditPage)))
	mux.Post("/admin/edit", dynamicMiddleware.Then(http.HandlerFunc(app.streamEditForm)))
	mux.Get("/admin/add", dynamicMiddleware.Then(http.HandlerFunc(app.streamAddPage)))
	mux.Post("/admin/add", dynamicMiddleware.Then(http.HandlerFunc(app.streamAddForm)))

	mux.Get("/static/", app.static())

	return app.Session.Enable(standardMiddleware.Then(mux))
}
