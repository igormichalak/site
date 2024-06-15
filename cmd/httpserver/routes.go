package main

import (
	"net/http"
	"os"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServerFS(os.DirFS("./public"))
	mux.Handle("GET /", fileServer)

	mux.HandleFunc("GET /{$}", app.homeView)
	mux.HandleFunc("GET /sitemap.xml", app.sitemap)
	mux.HandleFunc("GET /{slug}", app.postView)

	return app.recoverPanic(securityHeaders(mux))
}
