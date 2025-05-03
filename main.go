// main.go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := NewConfig("0.0.0.0", 45000)
	server := NewTimeServer(config)

	fmt.Println("=================================================")
	fmt.Println("       TIME SERVER - Concurrent TCP Server       ")
	fmt.Println("=================================================")
	fmt.Printf("Listening on: %s\n", config.Address())
	fmt.Println("Commands supported:")
	fmt.Println("  - TIME    : Get current server time")
	fmt.Println("  - QUIT    : Close the connection")
	fmt.Println("=================================================")

	go func() {
		if err := server.Start(); err != nil {
			fmt.Printf("[MAIN] Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("Server is running. Press Ctrl+C to stop...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\n[MAIN] Shutdown signal received")
	server.Stop()
}
