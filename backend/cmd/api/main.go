package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"
	"strings"

	"github.com/AdityaVarmaUddaraju/paytm/internal/data"
	_ "github.com/lib/pq"
)

type config struct {
	port         int
	env          string
	jwtSecretKey string
	db           struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	cfg    config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.jwtSecretKey, "jwt-secret-key", "", "JWT secret key")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL dsn")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space seperated)", func(s string) error {
		cfg.cors.trustedOrigins = strings.Fields(s)
		return nil
	})

	flag.Parse()

	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(cfg)

	if err != nil {
		jsonLogger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	jsonLogger.Info("database connection established")

	app := &application{
		cfg:    cfg,
		logger: jsonLogger,
		models: data.NewModels(db),
	}

	err = app.server()

	if err != nil {
		jsonLogger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
