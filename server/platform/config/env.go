package config

import (
	"os"
	"strconv"
)

type EnvSettings struct {
	APIKey string
	Prod   bool
	Auth   bool
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
	prod, err := strconv.ParseBool(GetKey("PROD", "false"))
	if err != nil {
		panic(err)
	}
	auth, err := strconv.ParseBool(GetKey("AUTH", "true"))
	if err != nil {
		panic(err)
	}
	return EnvSettings{
		APIKey: MustGetKey("API_KEY"),
		Prod:   prod,
		Auth:   auth,
	}
}
