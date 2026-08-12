package config

type App struct {
	Settings EnvSettings
}

func NewApp() *App {
	return &App{
		Settings: newEnvSettings(),
	}
}
