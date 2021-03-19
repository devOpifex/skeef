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

func (app *Application) getAuthenticated(r *http.Request) string {
	auth := app.Session.Get(r, "authenticatedUserID")
	return fmt.Sprintf("%v", auth)
}

func (app *Application) adminPage(w http.ResponseWriter, r *http.Request) {
	if !app.isAuthenticated(r) {
		http.Redirect(w, r, "/admin/signin", http.StatusSeeOther)
		return
	}

	email := app.getAuthenticated(r)

	if app.License.Email == "" {
		license, err := app.Database.GetLicense(email)

		if err != nil {
			http.Error(w, "Could not fetch license", http.StatusInternalServerError)
			return
		}

		app.License = license
	}

	tmplData := templateData{}
	tmplData.License = app.License

	app.render(w, r, []string{"ui/html/admin.page.tmpl"}, tmplData)
}
