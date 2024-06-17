package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var untrackedPaths = []string{"/search", "/css", "/fonts", "/js", "/favicon.ico"}

func (app *application) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := []string{
			"default-src 'self'",
			"script-src 'self'",
			"frame-src 'self'",
			"style-src 'self' 'unsafe-inline'",
			"connect-src 'self'",
		}

		w.Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Host, Origin, Referer, Accept, Content-Type, User-Agent, Cookie, X-Csrf-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Add("Vary", "Origin")

		// Preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.error(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) wwwRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Debug && !strings.HasPrefix(r.Host, "www.") {
			dst := "https://www." + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, dst, http.StatusMovedPermanently)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func formatDuration(d time.Duration) string {
	if d.Seconds() >= 1 {
		return fmt.Sprintf("%.2fs", float64(d.Milliseconds())/1e3)
	}
	if d.Milliseconds() >= 1 {
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1e3)
	}
	return fmt.Sprintf("%dμs", d.Microseconds())
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range untrackedPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		start := time.Now()

		next.ServeHTTP(w, r)

		var (
			method  = r.Method
			uri     = r.URL.RequestURI()
			took    = formatDuration(time.Now().Sub(start))
			referer = r.Referer()
		)

		app.Logger.Info("Request", "method", method, "uri", uri, "referer", referer, "took", took)
	})
}
