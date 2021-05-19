package app

import (
	"net/http"
	"strconv"
	"strings"

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

	ok := networkTypesOk(
		checkboxToInt(r.Form.Get("retweetsNet")),
		checkboxToInt(r.Form.Get("mentionsNet")),
		checkboxToInt(r.Form.Get("hashtagsNet")),
	)

	if ok {
		err = app.Database.UpdateStream(
			r.Form.Get("track"),
			r.Form.Get("follow"),
			r.Form.Get("locations"),
			r.Form.Get("name"),
			r.Form.Get("currentName"),
			r.Form.Get("exclude"),
			r.Form.Get("desc"),
			maxEdges,
			checkboxToInt(r.Form.Get("retweetsNet")),
			checkboxToInt(r.Form.Get("mentionsNet")),
			checkboxToInt(r.Form.Get("hashtagsNet")),
			r.Form.Get("filterLevel"),
		)

		if err != nil {
			tmplData.Errors["failure"] = "Failed to update stream"
		} else {
			tmplData.Flash["success"] = "Successfully updated stream"
		}
	} else {
		tmplData.Errors["failure"] = "Must check at least one of retweets, mentions, or hashtags."
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

func checkboxToInt(value string) int {
	return len(value)
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
	desc := r.Form.Get("desc")
	retweetsNet := checkboxToInt(r.Form.Get("retweetsNet"))
	mentionsNet := checkboxToInt(r.Form.Get("mentionsNet"))
	hashtagsNet := checkboxToInt(r.Form.Get("hashtagsNet"))
	filterLevel := r.Form.Get("filterLevel")

	ok := networkTypesOk(retweetsNet, mentionsNet, hashtagsNet)

	if !ok {
		tmplData.Errors["stream"] = "Must check at least one of retweets, mentions, or hashtags"
	}

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
		err = app.Database.InsertStream(name, follow, track, locations, exclude, maxEdges, desc, retweetsNet, mentionsNet, hashtagsNet, filterLevel)

		if err != nil {
			app.ErrorLog.Println(err)
			tmplData.Errors["stream"] = "Failed to add the stream to the database"
		} else {
			tmplData.Flash["stream"] = "Stream added to the database"
		}
	}

	app.render(w, r, []string{"ui/html/add.page.tmpl"}, tmplData)
}

func exclusionMap(exclusion string) map[string]bool {
	mp := make(map[string]bool)
	list := strings.Split(exclusion, ",")

	for _, str := range list {
		key := strings.TrimSpace(str)
		mp[key] = true
	}

	return mp
}

func networkTypesOk(retweetsNet, mentionsNet, hashtagsNet int) bool {
	total := retweetsNet + mentionsNet + hashtagsNet
	return total != 0
}
