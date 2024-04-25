package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	shutdownErr := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		s := <-quit

		app.logger.Info(
			"shutting down the server",
			"signal", s.String(),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)

		if err != nil {
			shutdownErr <- err
		}

		app.logger.Info(
			"completing background tasks",
			"addr", srv.Addr,
		)

		app.wg.Wait()

		shutdownErr <- nil
	}()

	app.logger.Info(
		"starting server",
		"addr", srv.Addr,
		"env", app.cfg.env,
	)

	err :=  srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr

	if err != nil {
		return err
	}

	app.logger.Info(
		"server stopped",
		"addr", srv.Addr,
	)

	return nil
}