package main

import (
	"net/http"

	"github.com/igormichalak/site/view"
)

func (app *application) homeView(w http.ResponseWriter, r *http.Request) {
	if err := view.Home().Render(r.Context(), w); err != nil {
		app.error(w, r, err)
	}
}

func (app *application) sitemap(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("sitemap"))
}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	if slug == "" {
		http.NotFound(w, r)
		return
	}

	entry, ok := view.PostIndex[slug]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if err := view.Post(entry).Render(r.Context(), w); err != nil {
		app.error(w, r, err)
	}
}
