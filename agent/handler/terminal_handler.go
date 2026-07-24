package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// TerminalHandler handles WebSocket terminal connections.
type TerminalHandler struct {
	sessions   map[string]*TerminalSession
	sessionsMu sync.RWMutex
}

// TerminalSession represents an active terminal session.
type TerminalSession struct {
	ID        string
	CreatedAt time.Time
	conn      *websocket.Conn
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	done      chan struct{}
}

// NewTerminalHandler creates a new TerminalHandler.
func NewTerminalHandler() *TerminalHandler {
	return &TerminalHandler{
		sessions: make(map[string]*TerminalSession),
	}
}

// HandleTerminal handles WebSocket terminal connections.
func (h *TerminalHandler) HandleTerminal(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	sessionID := c.Query("session_id")
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	session := &TerminalSession{
		ID:        sessionID,
		CreatedAt: time.Now(),
		conn:      conn,
		done:      make(chan struct{}),
	}

	h.sessionsMu.Lock()
	h.sessions[sessionID] = session
	h.sessionsMu.Unlock()

	defer func() {
		conn.Close()
		h.sessionsMu.Lock()
		delete(h.sessions, sessionID)
		h.sessionsMu.Unlock()
		if session.cmd != nil && session.cmd.Process != nil {
			session.cmd.Process.Kill()
		}
	}()

	// Start a shell process
	cmd := exec.Command("/bin/bash", "-i")
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("Failed to create stdin pipe: %v", err)
		return
	}
	session.stdin = stdinPipe

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to create stdout pipe: %v", err)
		return
	}
	session.stdout = stdoutPipe

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Failed to create stderr pipe: %v", err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start shell: %v", err)
		return
	}
	session.cmd = cmd

	// Goroutine to read from stdout and send to WebSocket
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stdoutPipe.Read(buf)
			if n > 0 {
				conn.WriteMessage(websocket.TextMessage, buf[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	// Goroutine to read from stderr and send to WebSocket
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stderrPipe.Read(buf)
			if n > 0 {
				conn.WriteMessage(websocket.TextMessage, buf[:n])
			}
			if err != nil {
				break
			}
		}
	}()

	// Read from WebSocket and write to stdin
	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// Handle resize commands
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err == nil {
				if op, ok := msg["op"].(string); ok && op == "resize" {
					if _, ok := msg["rows"].(float64); ok {
						if _, ok := msg["cols"].(float64); ok {
							// Resize terminal (would need pty for real resize)
						}
					}
					continue
				}
			}

			stdinPipe.Write(message)
		}
	}()

	// Wait for command to finish
	cmd.Wait()
	close(session.done)
}

// GetSessions returns a list of active terminal sessions.
func (h *TerminalHandler) GetSessions(c *gin.Context) {
	h.sessionsMu.RLock()
	defer h.sessionsMu.RUnlock()

	sessions := make([]map[string]interface{}, 0, len(h.sessions))
	for _, s := range h.sessions {
		sessions = append(sessions, map[string]interface{}{
			"id":         s.ID,
			"created_at": s.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": sessions})
}

// MetricsWSHandler handles WebSocket metrics streaming.
type MetricsWSHandler struct {
	metricsService interface {
		GetOverview() (interface{}, error)
	}
}

// NewMetricsWSHandler creates a new MetricsWSHandler.
func NewMetricsWSHandler(metricsService interface{}) *MetricsWSHandler {
	return &MetricsWSHandler{metricsService: metricsService}
}

// HandleMetricsWS handles WebSocket metrics streaming.
func (h *MetricsWSHandler) HandleMetricsWS(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			data := map[string]interface{}{
				"timestamp": time.Now().Unix(),
				"type":      "metrics",
			}

			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}
