package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.cfg.port),
		ErrorLog: slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, nil), slog.LevelError),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Info(
		"starting server",
		"addr", srv.Addr,
		"env", app.cfg.env,
	)

	return srv.ListenAndServe()
}