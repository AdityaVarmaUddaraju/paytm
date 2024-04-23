package main

import (
	"flag"
	"log/slog"
	"os"
)

type config struct {
	port int
	env  string
}

type application struct  {
	cfg config
	logger *slog.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	jsonLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &application{
		cfg: cfg,
		logger: jsonLogger,
	}

	err := app.server()

	if err != nil {
		jsonLogger.Error(err.Error())
		os.Exit(1)
	}	
}
