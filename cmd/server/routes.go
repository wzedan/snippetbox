package main

import "net/http"

func (app *Application) buildMux() *http.ServeMux {
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/create", app.snippetCreate)
	mux.HandleFunc("/view", app.snippetView)
	return mux
}
