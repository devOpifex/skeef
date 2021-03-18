package app

import (
	"net/http"
	"regexp"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (app *Application) setupPage(w http.ResponseWriter, r *http.Request) {

	if app.Setup.Tables {
		http.Redirect(w, r, "/setup/validate", http.StatusSeeOther)
		return
	}

	errUser := app.Database.CreateTableUser()
	errLicense := app.Database.CreateTableLicense()

	if errUser != nil || errLicense != nil {
		http.Error(w, "Failed to create database tables", http.StatusInternalServerError)
		return
	}

	app.Setup.Tables = true

	app.render(w, r, []string{"ui/html/setup.page.tmpl"}, templateData{})
}

func (app *Application) setupForm(w http.ResponseWriter, r *http.Request) {

	if app.Setup.Admin {
		http.Error(w, "Admin user already created", http.StatusNotAcceptable)
		return
	}

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
		tmplData.Errors["password"] = "Password must be at least 5 characters"
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
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
	}

	app.Setup.Admin = true

	http.Redirect(w, r, "/setup/validate", http.StatusSeeOther)
}

func (app *Application) validatePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, []string{"ui/html/validate.page.tmpl"}, templateData{})
}

func (app *Application) validateForm(w http.ResponseWriter, r *http.Request) {

	if app.Setup.License {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusInternalServerError)
	}

	email := r.PostForm.Get("email")
	license := r.PostForm.Get("license")

	app.Database.InsertLicense(email, license)
	app.Setup.License = true

	app.Session.Put(r, "authenticatedUserID", email)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
