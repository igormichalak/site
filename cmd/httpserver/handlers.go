package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home"))
}

func (app *application) sitemap(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sitemap"))
}
