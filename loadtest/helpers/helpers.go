package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func pingForCompile(url string, source []byte, apiKey string) (int, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(source))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if _, err := io.Copy(io.Discard, res.Body); err != nil {
		return 0, err
	}

	return res.StatusCode, nil
}

func TestCompile(url string, apiKey string, concurrency int, totalRequests int) {
	fmt.Println("URL:", url)
	fmt.Println("Concurrency:", concurrency)
	fmt.Println("Total requests:", totalRequests)

	texPath := "loadtest/resume.tex"
	source, err := os.ReadFile(texPath)
	if err != nil {
		fmt.Println("failed to read resume.tex:", err)
		return
	}

	body, err := json.Marshal(map[string]string{"source": string(source)})
	if err != nil {
		fmt.Println("failed to marshal request body:", err)
		return
	}

	var latencySum time.Duration
	var successes int
	var wg sync.WaitGroup
	var mu sync.Mutex

	sem := make(chan struct{}, concurrency)
	start := time.Now()

	for range totalRequests {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()

			reqStart := time.Now()
			statusCode, err := pingForCompile(url, body, apiKey)
			if err != nil || statusCode >= 400 {
				if err != nil {
					fmt.Println("error:", err)
				} else {
					fmt.Println("bad status:", statusCode)
				}
				return
			}
			reqElapsed := time.Since(reqStart)
			mu.Lock()
			latencySum += reqElapsed
			successes++
			mu.Unlock()
		})
	}

	wg.Wait()

	elapsed := time.Since(start)
	failures := totalRequests - successes
	rps := float64(totalRequests) / elapsed.Seconds()
	avgLatency := latencySum / time.Duration(totalRequests)

	fmt.Println("Time elapsed:", elapsed)
	fmt.Println("Successes:", successes, "|", "Failures", failures)
	fmt.Printf("Requests/sec: %.2f\n", rps)
	fmt.Printf("Requests/min: %.2f\n", rps*60)
	fmt.Printf("Average latency: %v/req\n", avgLatency)
}
