package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devOpifex/opiflex/config"
	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/sessions"
)

const (
	sessionName     = "example-twtter-app"
	sessionSecret   = "example cookie signing secret"
	sessionUserKey  = "twitterID"
	sessionUsername = "twitterUsername"
)

// Application Application object
type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Config   config.Config
	Oauth    oauth1.Config
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world")
}

func (app *Application) profile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are logged in")
}

// Handlers Returns all routes
func (app *Application) Handlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/profile", app.profile)
	mux.Handle("/login", twitter.LoginHandler(&app.Oauth, nil))
	mux.Handle("/"+app.Config.TwitterCallbackPath, twitter.CallbackHandler(&app.Oauth, issueSession(), nil))
	return mux
}

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		twitterUser, err := twitter.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = twitterUser.ID
		session.Values[sessionUsername] = twitterUser.ScreenName
		session.Save(w)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
