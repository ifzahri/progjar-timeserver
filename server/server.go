package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"
)

// sentinel errors
var (
	ErrServerAlreadyRunning = errors.New("server already running")
	ErrServerNotRunning     = errors.New("server not running")
)

// TimeServer
type TimeServer struct {
	config      Config
	status      bool
	wg          sync.WaitGroup
	mutex       sync.RWMutex
	server      net.Listener
	logger      Logger
	clientID    int
	ctx         context.Context
	cancelFunc  context.CancelFunc
	idleTimeout time.Duration
}

// NewTimeServer
func NewTimeServer(cfg Config) *TimeServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &TimeServer{
		config:      cfg,
		logger:      *NewLogger("SERVER"),
		clientID:    0,
		ctx:         ctx,
		cancelFunc:  cancel,
		idleTimeout: 30 * time.Second,
	}
}

// Start server
func (ts *TimeServer) Start() error {
	ts.mutex.Lock()
	if ts.status {
		ts.mutex.Unlock()
		return ErrServerAlreadyRunning
	}
	ts.mutex.Unlock()

	listener, err := net.Listen("tcp", ts.config.Address())
	if err != nil {
		return fmt.Errorf("could not start server: %w", err)
	}

	ts.mutex.Lock()
	ts.server = listener
	ts.status = true
	ts.mutex.Unlock()

	ts.logger.Info("Time server started on %s", ts.config.Address())

	go ts.ConnectionHandler()
	return nil
}

// Stop server
func (ts *TimeServer) Stop() error {
	ts.mutex.Lock()
	if !ts.status {
		ts.mutex.Unlock()
		return ErrServerNotRunning
	}

	ts.status = false
	ts.cancelFunc() // cancel connection

	if ts.server != nil {
		ts.logger.Info("Closing server listener")
		err := ts.server.Close()
		if err != nil {
			ts.mutex.Unlock()
			return fmt.Errorf("error closing listener: %w", err)
		}
	}
	ts.mutex.Unlock()

	// wait goroutine
	c := make(chan struct{})
	go func() {
		ts.wg.Wait()
		close(c)
	}()

	// timeout
	select {
	case <-c:
		// all connection done
	case <-time.After(5 * time.Second):
		ts.logger.Warn("Forced shutdown after timeout waiting for connections")
	}

	ts.logger.Info("Server stopped gracefully")
	return nil
}

// CheckServer
func (ts *TimeServer) CheckServer() bool {
	ts.mutex.RLock() // Read only lock
	defer ts.mutex.RUnlock()
	return ts.status
}

// IdleTimeout
func (ts *TimeServer) IdleTimeout(duration time.Duration) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	ts.idleTimeout = duration
}

// ConnectionHandler
func (ts *TimeServer) ConnectionHandler() {
	for {
		conn, err := ts.server.Accept()
		if err != nil {
			// if server not running
			if !ts.CheckServer() {
				break
			}
			ts.logger.Error("Error accepting connection: %v", err)
			// backoff
			time.Sleep(100 * time.Millisecond)
			continue
		}

		ts.mutex.Lock()
		clientID := ts.clientID
		ts.clientID++
		ts.mutex.Unlock()

		ts.wg.Add(1)
		go ts.ClientHandler(conn, clientID)
	}
}

// ClientHandler
func (ts *TimeServer) ClientHandler(conn net.Conn, clientID int) {
	defer ts.wg.Done()
	defer conn.Close()

	clientLogger := NewLogger(fmt.Sprintf("CLIENT-%d", clientID))
	clientAddr := conn.RemoteAddr().String()
	clientLogger.Info("New connection from %s", clientAddr)

	// Banner
	welcomeMsg := "Welcome to Time Server. Send 'TIME' for current time, 'QUIT' to disconnect.\r\n"
	if _, err := conn.Write([]byte(welcomeMsg)); err != nil {
		clientLogger.Error("Error sending welcome message: %v", err)
		return
	}

	buffer := make([]byte, 1024)

	// context to stop goroutine
	for {
		select {
		case <-ts.ctx.Done():
			clientLogger.Info("Connection closed due to server shutdown")
			return
		default:
			// deadline
			if err := conn.SetReadDeadline(time.Now().Add(ts.idleTimeout)); err != nil {
				clientLogger.Error("Failed to set read deadline: %v", err)
				return
			}

			n, err := conn.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// timeout, then continue loop
					continue
				}
				if err != io.EOF {
					clientLogger.Warn("Error reading from client: %v", err)
				}
				return
			}

			if n == 0 {
				continue
			}

			command := strings.ToUpper(strings.TrimSpace(string(buffer[:n])))
			clientLogger.Debug("Received command: %s", command)

			if strings.Contains(command, "TIME") {
				ts.TimeHandler(conn, *clientLogger)
			}
			if strings.Contains(command, "QUIT") {
				clientLogger.Info("Client requested disconnect")
				return
			}
		}
	}
}

// TimeHandler
func (ts *TimeServer) TimeHandler(conn net.Conn, logger Logger) {
	currentTime := time.Now().Format("15:04:05")
	response := fmt.Sprintf("JAM %s\r\n", currentTime)

	if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		logger.Error("Failed to set write deadline: %v", err)
		return
	}

	_, err := conn.Write([]byte(response))
	if err != nil {
		logger.Error("Error sending time: %v", err)
	} else {
		logger.Info("Time sent successfully: JAM %s", currentTime)
	}
}
