package app

import "net/http"

func (app *Application) static() http.Handler {
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	return http.StripPrefix("/static", fileServer)
}
