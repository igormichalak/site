package main

import (
	"net/http"

	"github.com/igormichalak/site/view"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if err := view.Home().Render(r.Context(), w); err != nil {
		app.error(w, r, err)
	}
}

func (app *application) sitemap(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sitemap"))
}
