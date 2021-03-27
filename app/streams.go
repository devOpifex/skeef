package app

import (
	"net/http"
	"strconv"

	"github.com/devOpifex/skeef-app/stream"
)

func (app *Application) streamEditPage(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) streamEditForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var tmplData templateData
	tmplData.Errors = make(map[string]string)
	tmplData.Flash = make(map[string]string)

	maxEdges, _ := strconv.Atoi(r.Form.Get("maxEdges"))

	err = app.Database.UpdateStream(
		r.Form.Get("track"),
		r.Form.Get("follow"),
		r.Form.Get("locations"),
		r.Form.Get("name"),
		r.Form.Get("currentName"),
		maxEdges,
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
		MaxEdges:  maxEdges,
	}

	app.render(w, r, []string{"ui/html/stream.page.tmpl"}, tmplData)
	http.Redirect(w, r, "/admin/edit/"+r.Form.Get("name"), http.StatusSeeOther)
}

func (app *Application) streamAddPage(w http.ResponseWriter, r *http.Request) {
	if !app.isAuthenticated(r) {
		http.Redirect(w, r, "/admin/signin", http.StatusSeeOther)
		return
	}

	app.render(w, r, []string{"ui/html/add.page.tmpl"}, templateData{})
}

func (app *Application) streamAddForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var tmplData templateData
	tmplData.Errors = make(map[string]string)
	tmplData.Flash = make(map[string]string)

	name := r.Form.Get("name")
	follow := r.Form.Get("follow")
	track := r.Form.Get("track")
	locations := r.Form.Get("locations")
	exclude := r.Form.Get("exclude")
	maxEdges := r.Form.Get("maxEdges")

	if name == "" {
		tmplData.Errors["stream"] = "Must specify a name"
	}

	if follow == "" && track == "" {
		tmplData.Errors["stream"] = "Must use 'follow' or 'track' (or both)"
	}

	streamExists, err := app.Database.StreamExists(name)

	if err != nil {
		tmplData.Errors["stream"] = "Failed to check if stream exists"
	}

	if streamExists {
		tmplData.Errors["stream"] = "A stream under that name already exists"
	}

	if len(tmplData.Errors) == 0 {
		err = app.Database.InsertStream(name, follow, track, locations, exclude, maxEdges)

		if err != nil {
			tmplData.Errors["stream"] = "Failed to add the stream to the database"
		} else {
			tmplData.Flash["stream"] = "Stream added to the database"
		}
	}

	app.render(w, r, []string{"ui/html/add.page.tmpl"}, tmplData)
}
