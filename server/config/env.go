package config

import (
	"os"
)

type EnvSettings struct {
	APIKey string
}

func MustGetKey(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("missing required env var:" + key)
	}
	return val
}

func GetKey(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func newEnvSettings() EnvSettings {
	return EnvSettings{
		APIKey: MustGetKey("API_KEY"),
	}
}
