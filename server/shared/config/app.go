package config

import (
	"log/slog"
	"os"
)

type App struct {
	Settings EnvSettings
}

func setupLogger(prod bool) {
	var logger slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if prod {
		logger = slog.NewJSONHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(logger))
}

func NewApp() *App {
	app := App{
		Settings: newEnvSettings(),
	}

	setupLogger(app.Settings.Prod)

	return &app
}
