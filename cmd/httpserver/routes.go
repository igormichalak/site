package main

import (
	"net/http"
	"os"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServerFS(os.DirFS("./public"))
	mux.Handle("GET /", fileServer)

	return app.recoverPanic(securityHeaders(mux))
}
