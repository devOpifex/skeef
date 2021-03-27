package app

import (
	"fmt"
	"net/http"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

func (app *Application) signinPage(w http.ResponseWriter, r *http.Request) {
	tmpls := []string{
		"ui/html/signin.page.tmpl",
	}

	app.render(w, r, tmpls, templateData{})
}

func (app *Application) signinForm(w http.ResponseWriter, r *http.Request) {
	tmpls := []string{
		"ui/html/signin.page.tmpl",
	}

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var tmplData templateData
	tmplData.Errors = make(map[string]string)

	email, err := app.Database.Authenticate(r.PostForm.Get("email"), r.PostForm.Get("password"))

	if err != nil {
		tmplData.Errors["credentials"] = "Invalid credentials"
		app.render(w, r, tmpls, tmplData)
		return
	}

	app.Session.Put(r, "authenticatedUserID", email)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r, "authenticatedUserID")
}

// will eventually be useful
func (app *Application) GetAuthenticated(r *http.Request) string {
	auth := app.Session.Get(r, "authenticatedUserID")
	return fmt.Sprintf("%v", auth)
}

func (app *Application) adminPage(w http.ResponseWriter, r *http.Request) {
	if !app.isAuthenticated(r) {
		http.Redirect(w, r, "/admin/signin", http.StatusSeeOther)
		return
	}

	if app.License.Email == "" {
		license, err := app.Database.GetLicense()

		if err != nil {
			http.Error(w, "Could not fetch license", http.StatusInternalServerError)
			return
		}

		app.License = license
	}

	hasTokens := app.Database.TokensExist()
	tmplData := templateData{}
	tmplData.License = app.License
	tmplData.HasTokens = hasTokens
	tmplData.Email = app.GetAuthenticated(r)

	if hasTokens {
		streams, err := app.Database.GetStreams()

		if err != nil {
			tmplData.Errors["existingStreams"] = "Could not retrieve streams"
		}

		tmplData.Streams = streams
	}

	app.render(w, r, []string{"ui/html/admin.page.tmpl"}, tmplData)
}

func (app *Application) adminForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var tmplData templateData
	tmplData.Errors = make(map[string]string)
	tmplData.Flash = make(map[string]string)
	tmplData.Email = app.GetAuthenticated(r)

	hasTokens := app.Database.TokensExist()

	action := r.Form.Get("action")
	if action == "twitter" {
		apiKey := r.Form.Get("apiKey")
		apiSecret := r.Form.Get("apiSecret")
		accessToken := r.Form.Get("accessToken")
		accessSecret := r.Form.Get("accessSecret")

		// UPDATE OR INSERT TOKENS
		if hasTokens {
			err = app.Database.UpdateTokens(apiKey, apiSecret, accessToken, accessSecret)

			if err != nil {
				app.ErrorLog.Println(err)
				tmplData.Errors["any"] = "Could not store data"
			}
		} else {
			err = app.Database.InsertTokens(apiKey, apiSecret, accessToken, accessSecret)

			if err != nil {
				app.ErrorLog.Println(err)
				tmplData.Errors["any"] = "Could not store data"
			}
		}
	}

	if action == "license" {
		newLicense := r.Form.Get("license")

		oldLicense := app.License.License
		app.License.License = newLicense
		response := app.LicenseCheck(false)

		if !response.Success {
			tmplData.Errors["license"] = response.Reason
			app.License.License = oldLicense
		} else {
			err = app.Database.UpdateLicense(app.GetAuthenticated(r), newLicense)

			if err != nil {
				tmplData.Errors["license"] = "Failed to update license"
			}
		}

	}

	if action == "validity" {
		response := app.LicenseCheck(false)
		app.LicenseResponse = response

		tmplData.Flash["validity"] = response.Reason
	}

	if action == "newPassword" {
		password := r.PostForm.Get("password")
		password2 := r.PostForm.Get("password2")

		if password == "" || password2 == "" {
			tmplData.Errors["password"] = "Empty password"
		}

		if password != password2 {
			tmplData.Errors["password"] = "Passwords do not match"
		}

		if utf8.RuneCountInString(password) < 5 {
			tmplData.Errors["password"] = "Password must be at least 5 characters"
		}

		if len(tmplData.Errors) > 0 {
			app.render(w, r, []string{"ui/html/profile.page.tmpl"}, tmplData)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			http.Error(w, "Could not hash password", http.StatusInternalServerError)
			return
		}

		err = app.Database.ChangePassword(app.License.Email, string(hashedPassword))

		if err != nil {
			http.Error(w, "Could not change password", http.StatusInternalServerError)
			return
		}

		tmplData.Flash["password"] = "Password changed!"
	}

	if action == "deleteStream" {
		err = app.Database.DeleteStream(r.Form.Get("streamName"))

		if err != nil {
			tmplData.Errors["existingStreams"] = "Failed to delete the stream from the database"
		} else {
			tmplData.Flash["existingStreams"] = "Deleted stream from the database"
		}
	}

	if action == "startStream" {

		if app.Database.StreamOnGoing() {

			tmplData.Errors["existingStreams"] = "There is already one stream active"

		} else {
			err = app.Database.StartStream(r.Form.Get("streamName"))

			if err != nil {
				tmplData.Errors["existingStreams"] = "Failed to start stream"
			} else {

				app.LicenseValidity()
				app.Streaming = true

				tmplData.Flash["existingStreams"] = "Stream Started"
			}

		}

	}

	if action == "stopStream" {

		if !app.Database.StreamOnGoing() {
			tmplData.Errors["existingStreams"] = "There is no active stream to pause"
		} else {
			err = app.Database.PauseStream(r.Form.Get("streamName"))

			if err != nil {
				tmplData.Errors["existingStreams"] = "Failed to pause stream"
			} else {
				app.Streaming = false
				tmplData.Flash["existingStreams"] = "Stream Paused"
			}

		}

	}

	// get stored streams
	streams, err := app.Database.GetStreams()
	if err != nil {
		tmplData.Errors["existingStreams"] = "Could not retrieve streams"
	}
	tmplData.Streams = streams

	tmplData.License = app.License
	tmplData.HasTokens = app.Database.TokensExist()

	app.render(w, r, []string{"ui/html/admin.page.tmpl"}, tmplData)
}

func (app *Application) signout(w http.ResponseWriter, r *http.Request) {
	app.Session.Remove(r, "authenticatedUserID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
