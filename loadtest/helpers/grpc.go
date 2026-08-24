package helpers

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	compilerpb "jakes-resume-compiler/proto"
	"jakes-resume-compiler/server/shared/config"
)

func newClientConn(target string, useTLS bool) (*grpc.ClientConn, error) {
	if useTLS {
		creds := credentials.NewTLS(&tls.Config{})
		return grpc.NewClient(target, grpc.WithTransportCredentials(creds))
	}
	return grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// newClientConns opens count connections so requests are spread across separate
// HTTP/2 streams instead of being multiplexed over a single connection.
func newClientConns(target string, useTLS bool, count int) ([]*grpc.ClientConn, error) {
	conns := make([]*grpc.ClientConn, 0, count)
	for range count {
		conn, err := newClientConn(target, useTLS)
		if err != nil {
			for _, c := range conns {
				c.Close()
			}
			return nil, err
		}
		conns = append(conns, conn)
	}
	return conns, nil
}

func compileRPC(ctx context.Context, client compilerpb.CompilerClient, source string, apiKey string) error {
	ctx = metadata.AppendToOutgoingContext(ctx, config.APIKeyHeader, apiKey)

	_, err := client.Compile(ctx, &compilerpb.CompileRequest{TexSource: source})
	return err
}

func TestGRPCCompile(target string, apiKey string, concurrency int, totalRequests int, useTLS bool) {
	fmt.Println("Target:", target)
	fmt.Println("Concurrency:", concurrency)
	fmt.Println("Total requests:", totalRequests)

	texPath := "loadtest/resume.tex"
	source, err := os.ReadFile(texPath)
	if err != nil {
		fmt.Println("failed to read resume.tex:", err)
		return
	}

	conns, err := newClientConns(target, useTLS, concurrency)

	if err != nil {
		fmt.Println("failed to create clients:", err)
		return
	}
	defer func() {
		for _, conn := range conns {
			conn.Close()
		}
	}()

	clients := make([]compilerpb.CompilerClient, len(conns))
	for i, conn := range conns {
		clients[i] = compilerpb.NewCompilerClient(conn)
	}

	var latencySum time.Duration
	var successes int
	var wg sync.WaitGroup
	var mu sync.Mutex

	sem := make(chan struct{}, concurrency)
	start := time.Now()

	for i := range totalRequests {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()

			client := clients[i%len(clients)]

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
