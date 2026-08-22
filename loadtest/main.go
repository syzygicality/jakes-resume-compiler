package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"jakes-resume-compiler/loadtest/helpers"
)

// run with: go run ./loadtest

type TestConfig struct {
	ConcurrencyTests []int  `json:"concurrencyTests"`
	TotalRequests    int    `json:"totalRequests"`
	Host             string `json:"host"`
}

func main() {
	var cfg TestConfig
	data, err := os.ReadFile("loadtest/test-config.json")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	mode := os.Getenv("MODE")
	if mode == "HTTP" {
		url := fmt.Sprintf("%s/compile", cfg.Host)

		for _, c := range cfg.ConcurrencyTests {
			helpers.TestHTTPCompile(url, apiKey, c, cfg.TotalRequests)
		}
		return
	}

	target := strings.TrimPrefix(strings.TrimPrefix(cfg.Host, "http://"), "https://")

	for _, c := range cfg.ConcurrencyTests {
		helpers.TestGRPCCompile(target, apiKey, c, cfg.TotalRequests)
	}
}
