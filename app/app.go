package app

import (
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/devOpifex/skeef/db"
	"github.com/devOpifex/skeef/graph"
	"github.com/devOpifex/skeef/stream"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/golangcollege/sessions"
	"github.com/justinas/alice"
)

// Application Application object
type Application struct {
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	Database     db.Database
	Session      *sessions.Session
	Addr         string
	Count        int64
	Stream       *twitter.Stream
	Valid        bool
	Pool         *Pool
	Quit         chan struct{}
	Streaming    bool
	Graph        graph.Graph
	Connected    int
	Trend        map[int64]int
	Exclusion    map[string]bool
	MaxEdges     int
	StreamActive stream.Stream
	NotStreaming string
	TweetsUsers  []tweetsUsers
}

type Setup struct {
	Tables bool
	Admin  bool
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {

	if !app.Database.AdminExists() {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	var tmplData templateData
	tmplData.Flash = make(map[string]string)

	tmplData.Authenticated = app.isAuthenticated(r)

	tmplData.Streaming = app.Database.StreamOnGoing()

	// default message
	msg := app.NotStreaming
	if msg == "" {
		msg = "<h1 class='uppercase text-center font-light text-3xl pb-8'>Currently not streaming</h1>"
	}
	tmplData.Flash["message"] = msg

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
