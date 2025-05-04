package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ifzahri/progjar-timeserver/server"
)

func main() {
	host := flag.String("host", "0.0.0.0", "Host address to bind")
	port := flag.Int("port", 45000, "Port to listen on")
	logLevel := flag.String("loglevel", "info", "Log level (debug, info, warn, error)")

	flag.Parse()

	// config
	cfg := server.NewConfig(*host, *port)
	cfg.SetLogLevel(*logLevel)

	// server
	timeServer := server.NewTimeServer(*cfg)
	Banner(cfg)

	// goroutine
	go func() {
		if err := timeServer.Start(); err != nil {
			fmt.Printf("[FATAL] Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	// signal
	ShutdownSignal(timeServer)
}

// Banner
func Banner(cfg *server.Config) {
	fmt.Println("=================================================")
	fmt.Println("       TIME SERVER - Concurrent TCP Server       ")
	fmt.Println("=================================================")
	fmt.Printf("Listening on: %s\n", cfg.Address())
	fmt.Println("Commands supported:")
	fmt.Println("  - TIME    : Get current server time")
	fmt.Println("  - QUIT    : Close the connection")
	fmt.Println("=================================================")
	fmt.Println("Server is running. Press Ctrl+C to stop...")
}

// ShutdownSignal
func ShutdownSignal(srv *server.TimeServer) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// nterrupt
	sig := <-sigChan
	fmt.Printf("\n[MAIN] Shutdown signal received: %v\n", sig)

	// timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		if err := srv.Stop(); err != nil {
			fmt.Printf("[MAIN] Error during shutdown: %v\n", err)
		}
		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		fmt.Println("[MAIN] Shutdown timeout exceeded, forcing exit")
	case <-done:
		fmt.Println("[MAIN] Shutdown completed successfully")
	}
}
