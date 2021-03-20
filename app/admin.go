package app

import (
	"fmt"
	"net/http"
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

	hasTokens := app.Database.TokensExist()
	tmplData.HasTokens = hasTokens

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

	app.render(w, r, []string{"ui/html/admin.page.tmpl"}, tmplData)
}
