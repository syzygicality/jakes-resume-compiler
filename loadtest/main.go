package main

import (
	"os"

	"github.com/joho/godotenv"

	"jakes-resume-compiler/loadtest/helpers"
)

// run with: go run ./loadtest

func testResumakr() {
	url := "http://localhost:8080/compile"

	apiKey := os.Getenv("API_KEY")
	concurrencyTests := []int{20, 40, 60, 80, 100}
	totalRequests := 200

	for _, c := range concurrencyTests {
		helpers.TestCompile(url, apiKey, c, totalRequests)
	}
}

func main() {
	godotenv.Load()
	testResumakr()
}

// ==== BASELINE (ORIGINAL OPTIMIZATION COMMIT) ====

// URL: http://localhost:8080/compile
// Concurrency: 20
// Total requests: 200
// Time elapsed: 3.198909708s
// Successes: 200 | Failures 0
// Requests/sec: 62.52
// Requests/min: 3751.28
// Average latency: 313.649262ms/req
// URL: http://localhost:8080/compile
// Concurrency: 40
// Total requests: 200
// Time elapsed: 2.888454375s
// Successes: 200 | Failures 0
// Requests/sec: 69.24
// Requests/min: 4154.47
// Average latency: 555.563673ms/req
// URL: http://localhost:8080/compile
// Concurrency: 60
// Total requests: 200
// Time elapsed: 2.944711792s
// Successes: 200 | Failures 0
// Requests/sec: 67.92
// Requests/min: 4075.10
// Average latency: 842.058135ms/req
// URL: http://localhost:8080/compile
// Concurrency: 80
// Total requests: 200
// Time elapsed: 2.97715275s
// Successes: 200 | Failures 0
// Requests/sec: 67.18
// Requests/min: 4030.70
// Average latency: 1.107508158s/req
// URL: http://localhost:8080/compile
// Concurrency: 100
// Total requests: 200
// Time elapsed: 2.966351167s
// Successes: 200 | Failures 0
// Requests/sec: 67.42
// Requests/min: 4045.37
// Average latency: 1.36443438s/req
