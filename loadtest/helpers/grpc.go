package helpers

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	compilerpb "jakes-resume-compiler/proto"
	"jakes-resume-compiler/server/shared/config"
)

func compileRPC(ctx context.Context, client compilerpb.CompilerClient, source string, apiKey string) error {
	ctx = metadata.AppendToOutgoingContext(ctx, config.APIKeyHeader, apiKey)

	_, err := client.Compile(ctx, &compilerpb.CompileRequest{TexSource: source})
	return err
}

func TestGRPCCompile(target string, apiKey string, concurrency int, totalRequests int) {
	fmt.Println("Target:", target)
	fmt.Println("Concurrency:", concurrency)
	fmt.Println("Total requests:", totalRequests)

	texPath := "loadtest/resume.tex"
	source, err := os.ReadFile(texPath)
	if err != nil {
		fmt.Println("failed to read resume.tex:", err)
		return
	}

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("failed to create client:", err)
		return
	}
	defer conn.Close()

	client := compilerpb.NewCompilerClient(conn)

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
			if err := compileRPC(context.Background(), client, string(source), apiKey); err != nil {
				fmt.Println("error:", err)
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
