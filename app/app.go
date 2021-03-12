package app

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/devOpifex/skeef-app/config"
)

// Application Application object
type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Config   config.Config
}

//go:embed ui/html
var embededTemplates embed.FS

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"ui/html/home.page.html",
		"ui/html/base.layout.html",
		"ui/html/footer.partial.html",
	}

	ts, err := template.ParseFS(embededTemplates, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Handlers Returns all routes
func (app *Application) Handlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.Handle("/static/", app.static())
	return mux
}
