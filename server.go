package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type TimeServer struct {
	config   *Config
	status   bool
	wg       sync.WaitGroup // WaitGroup goroutines
	mutex    sync.Mutex     // Mutex for thread safety
	server   net.Listener   // Listener connections
	logger   *Logger        // Logger
	clientID int
}

func NewTimeServer(config *Config) *TimeServer {
	return &TimeServer{
		config:   config,
		logger:   NewLogger("SERVER"),
		clientID: 0,
	}
}

func (ts *TimeServer) Start() error {
	listener, err := net.Listen("tcp", ts.config.Address())
	if err != nil {
		return fmt.Errorf("could not start server: %v", err)
	}

	// mutex, then set server and status, then unlock
	ts.mutex.Lock()
	ts.server = listener
	ts.status = true
	ts.mutex.Unlock()
	ts.logger.Info("Time server started on %s", ts.config.Address())

	// connection loop
	ts.Connections()
	return nil
}

func (ts *TimeServer) Stop() {
	ts.mutex.Lock()
	if ts.status {
		ts.status = false
		if ts.server != nil {
			ts.logger.Info("Closing server listener")
			ts.server.Close()
		}
	}
	ts.mutex.Unlock()

	// Wait for all client
	ts.wg.Wait()
	ts.logger.Info("Server stopped gracefully")
}

func (ts *TimeServer) IsRunning() bool {
	// check server status
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	return ts.status
}

// connection handler
func (ts *TimeServer) Connections() {
	for ts.IsRunning() {
		// Accept new connections
		conn, err := ts.server.Accept()
		if err != nil {
			if ts.IsRunning() {
				ts.logger.Error("Error accepting connection: %v", err)
			}
			break
		}

		// get next client
		// increment id
		ts.mutex.Lock()
		clientID := ts.clientID
		ts.clientID++
		ts.mutex.Unlock()

		// put each conection on new goroutine
		ts.wg.Add(1)
		go ts.Client(conn, clientID)
	}
}

// Client handler
func (ts *TimeServer) Client(conn net.Conn, clientID int) {
	// Defer the wg.Done() to ensure it is called when the function exits
	defer ts.wg.Done()
	defer conn.Close()

	clientLogger := NewLogger(fmt.Sprintf("CLIENT-%d", clientID))
	clientAddr := conn.RemoteAddr().String()
	clientLogger.Info("New connection from %s", clientAddr)

	// Send welcome message to client
	welcomeMsg := "Welcome to Time Server. Send 'TIME' for current time, 'QUIT' to disconnect.\r\n"
	conn.Write([]byte(welcomeMsg))

	buffer := make([]byte, 1024)

	for ts.IsRunning() {
		// Set read deadline to prevent blocked reads from hanging forever
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		n, err := conn.Read(buffer)
		if err != nil {
			// Check if the error is a timeout
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			clientLogger.Warn("Error reading from client: %v", err)
			return
		}

		if n == 0 {
			continue
		}

		// Process command
		command := strings.ToUpper(strings.TrimSpace(string(buffer[:n])))

		if strings.Contains(command, "TIME") {
			ts.Time(conn, clientLogger)
		}
		if strings.Contains(command, "QUIT") {
			clientLogger.Info("Client requested disconnect")
			return
		}
	}
}

// Time handler
func (ts *TimeServer) Time(conn net.Conn, logger *Logger) {
	currentTime := time.Now().Format("15:04:05")
	response := []byte("JAM " + currentTime + "\r\n")

	_, err := conn.Write(response)
	if err != nil {
		logger.Error("Error sending time: %v", err)
	} else {
		logger.Info("Time sent successfully: JAM %s", currentTime)
	}
}
