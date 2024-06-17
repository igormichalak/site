package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DefaultPort = "8080"

type application struct {
	Logger *slog.Logger
	Debug  bool
}

func main() {
	certFile := flag.String("cert-file", "./cert.pem", "certificate file path")
	keyFile := flag.String("key-file", "./key.pem", "key file path")
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = DefaultPort
	}

	var logger *slog.Logger
	if *debug {
		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
		logger = slog.New(handler)
	} else {
		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(handler)
	}

	app := application{
		Logger: logger,
		Debug:  *debug,
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           app.routes(),
		ReadTimeout:       6 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      12 * time.Second,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    8_192,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
			CurvePreferences:   []tls.CurveID{tls.X25519, tls.CurveP256},
			ClientSessionCache: tls.NewLRUClientSessionCache(128),
		},
	}

	stopC := make(chan os.Signal, 1)
	signal.Notify(stopC, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		logger.Info(fmt.Sprintf("Starting server on %s...", srv.Addr))

		err := srv.ListenAndServeTLS(*certFile, *keyFile)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server failed", "err", err)
			os.Exit(1)
		}
	}()

	<-stopC
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "err", err)
		os.Exit(1)
	}
	logger.Info("Server gracefully stopped.")
}
