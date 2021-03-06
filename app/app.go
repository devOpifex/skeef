package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devOpifex/opiflex/config"
	oauth1Login "github.com/dghubble/gologin/v2/oauth1"
	gologin "github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/sessions"
)

// Constants for token
const (
	sessionName     = "opiflex"
	sessionSecret   = "opifex.org"
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
	mux.Handle("/login", gologin.LoginHandler(&app.Oauth, nil))
	mux.Handle("/"+app.Config.TwitterCallbackPath, gologin.CallbackHandler(&app.Oauth, issueSession(), nil))
	return mux
}

var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		twitterUser, err := gologin.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accessToken, accessSecret, err := oauth1Login.AccessTokenFromContext(ctx)

		fmt.Println(accessToken)
		fmt.Println(accessSecret)
		if err != nil {
			log.Print(err)
		}

		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = twitterUser.ID
		session.Values[sessionUsername] = twitterUser.ScreenName
		session.Save(w)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
