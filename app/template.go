package app

import (
	"embed"
	"net/http"
	"text/template"

	"github.com/devOpifex/skeef-app/db"
	"github.com/devOpifex/skeef-app/stream"
	"github.com/justinas/nosurf"
)

type templateData struct {
	Errors        map[string]string
	Authenticated bool
	CSRFToken     string
	License       db.License
	HasTokens     bool
	Flash         string
	Stream        stream.Stream
}

//go:embed ui/html
var embededTemplates embed.FS

func (app *Application) render(w http.ResponseWriter, r *http.Request, files []string, data templateData) {

	tmpls := []string{
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	tmpls = append(files, tmpls...)

	ts, err := template.ParseFS(embededTemplates, tmpls...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data.CSRFToken = nosurf.Token(r)

	err = ts.Execute(w, data)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
