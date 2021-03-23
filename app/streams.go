package app

import (
	"net/http"

	"github.com/devOpifex/skeef-app/stream"
)

func (app *Application) streamPage(w http.ResponseWriter, r *http.Request) {
	if !app.isAuthenticated(r) {
		http.Redirect(w, r, "/admin/signin", http.StatusSeeOther)
		return
	}

	name := r.URL.Query().Get(":stream")
	stream, err := app.Database.GetStream(name)

	if err != nil {
		http.Error(w, "Failed to fetch stream", http.StatusInternalServerError)
		return
	}

	var tmplData templateData
	tmplData.Stream = stream

	app.render(w, r, []string{"ui/html/stream.page.tmpl"}, tmplData)
}

func (app *Application) streamForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var tmplData templateData
	tmplData.Errors = make(map[string]string)
	tmplData.Flash = make(map[string]string)

	err = app.Database.UpdateStream(
		r.Form.Get("track"),
		r.Form.Get("follow"),
		r.Form.Get("locations"),
		r.Form.Get("name"),
		r.Form.Get("currentName"),
	)

	if err != nil {
		tmplData.Errors["failure"] = "Failed to update stream"
	} else {
		tmplData.Flash["success"] = "Successfully updated stream"
	}

	tmplData.Stream = stream.Stream{
		Follow:    r.Form.Get("follow"),
		Track:     r.Form.Get("track"),
		Locations: r.Form.Get("locations"),
		Name:      r.Form.Get("name"),
	}

	app.render(w, r, []string{"ui/html/stream.page.tmpl"}, tmplData)
	http.Redirect(w, r, "/admin/"+r.Form.Get("name"), http.StatusSeeOther)
}
