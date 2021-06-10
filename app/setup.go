package app

import (
	"net/http"
	"regexp"
	"unicode/utf8"

	"github.com/devOpifex/skeef-app/db"
	"golang.org/x/crypto/bcrypt"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (app *Application) setupPage(w http.ResponseWriter, r *http.Request) {

	if app.Database.AdminExists() {
		http.Redirect(w, r, "/setup/validate", http.StatusSeeOther)
		return
	}

	app.render(w, r, []string{"ui/html/setup.page.tmpl"}, templateData{})
}

func (app *Application) setupForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	password2 := r.PostForm.Get("password2")

	var tmplData templateData
	tmplData.Errors = make(map[string]string)

	if email == "" {
		tmplData.Errors["exists"] = "Empty email"
	}

	if password == "" || password2 == "" {
		tmplData.Errors["password"] = "Empty password"
	}

	if password != password2 {
		tmplData.Errors["password"] = "Passwords do not match"
	}

	if utf8.RuneCountInString(password) < 5 {
		tmplData.Errors["password"] = "Password must be at least 5 characters long"
	}

	if !EmailRX.MatchString(email) {
		tmplData.Errors["exists"] = "Invalid email address"
	}

	if len(tmplData.Errors) > 0 {
		app.render(w, r, []string{"ui/html/setup.page.tmpl"}, tmplData)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	err = app.Database.InsertUser(email, string(hashedPassword), 1)

	if err != nil {
		http.Error(w, "Could not create the user", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/setup/validate", http.StatusSeeOther)
}

func (app *Application) validatePage(w http.ResponseWriter, r *http.Request) {
	if app.Database.LicenseExists() {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.render(w, r, []string{"ui/html/validate.page.tmpl"}, templateData{})
}

func (app *Application) validateForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusInternalServerError)
	}

	tmplData := templateData{}
	tmplData.Errors = make(map[string]string)
	var response LicenseResponse

	email := r.PostForm.Get("email")
	license := r.PostForm.Get("license")

	if email == "" || license == "" {
		tmplData.Errors["exists"] = "Empty email or license"
	} else {
		app.License = db.License{
			Email:   email,
			License: license,
		}

		response = app.LicenseCheck(false)
	}

	if !response.Success {
		tmplData.Errors["license"] = response.Reason
	}

	if len(tmplData.Errors) > 0 {
		app.render(w, r, []string{"ui/html/validate.page.tmpl"}, tmplData)
		return
	}

	err = app.Database.InsertLicense(email, license)

	if err != nil {
		http.Error(w, "Could not store license", http.StatusInternalServerError)
		return
	}

	app.License.Email = email
	app.License.Email = license

	app.Session.Put(r, "authenticatedUserID", email)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
