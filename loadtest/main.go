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
	ConcurrencyTests []int    `json:"concurrencyTests"`
	TotalRequests    int      `json:"totalRequests"`
	Host             string   `json:"host"`
	Protocols        []string `json:"protocols"`
	UseTLS           bool     `json:"useTLS"`
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

	// The server speaks both protocols on one port, so an empty list means
	// exercise both against the same host.
	protocols := cfg.Protocols
	if len(protocols) == 0 {
		protocols = []string{"HTTP", "gRPC"}
	}

	for _, protocol := range protocols {
		fmt.Println("=====", protocol, "=====")

		switch protocol {
		case "HTTP":
			url := fmt.Sprintf("%s/compile", cfg.Host)

			for _, c := range cfg.ConcurrencyTests {
				helpers.TestHTTPCompile(url, apiKey, c, cfg.TotalRequests)
			}
		case "gRPC":
			target := strings.TrimPrefix(strings.TrimPrefix(cfg.Host, "http://"), "https://")

			for _, c := range cfg.ConcurrencyTests {
				helpers.TestGRPCCompile(target, apiKey, c, cfg.TotalRequests, cfg.UseTLS)
			}
		default:
			panic("protocols entries must be either HTTP or gRPC.")
		}
	}
}
