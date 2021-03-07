package app

import (
	"net/http"

	oauth1Login "github.com/dghubble/gologin/v2/oauth1"
	gologin "github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/sessions"
)

// Constants for session
const (
	sessionName         = "opiflex"
	sessionSecret       = "opifex.org"
	sessionUserKey      = "twitterID"
	sessionUsername     = "twitterUsername"
	twitterAccessToken  = "accessToken"
	twitterAccessSecret = "accessSecret"
)

// Session
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

func (app *Application) authenticate() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		twitterUser, err := gologin.UserFromContext(ctx)
		if err != nil {
			app.ErrorLog.Println("Failed to get user details")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accessToken, accessSecret, err := oauth1Login.AccessTokenFromContext(ctx)

		if err != nil {
			app.ErrorLog.Println("Failed to get access creds")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session := sessionStore.New(sessionName)
		session.Values[twitterAccessToken] = accessToken
		session.Values[twitterAccessSecret] = accessSecret
		session.Values[sessionUserKey] = twitterUser.ID
		session.Values[sessionUsername] = twitterUser.ScreenName
		sessionStore.Save(w, session)
		http.Redirect(w, r, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
