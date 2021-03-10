package app

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed ui/static
var embededFiles embed.FS

func (app *Application) static() http.Handler {
	fsys, err := fs.Sub(embededFiles, "static")

	if err != nil {
		app.ErrorLog.Panic("Internal error code: e-1")
		return nil
	}

	return http.FileServer(http.FS(fsys))
}
