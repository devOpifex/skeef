package app

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed ui/static
var embededStatic embed.FS

func (app *Application) static() http.Handler {
	fsys, err := fs.Sub(embededStatic, "ui/static")

	if err != nil {
		app.ErrorLog.Panic("Internal error code: e-1")
		return nil
	}

	return http.StripPrefix("/static", http.FileServer(http.FS(fsys)))
}
