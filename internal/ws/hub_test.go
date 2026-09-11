package ws_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sawitos/sawit/internal/ws"
)

func TestWebSocketHubConnection(t *testing.T) {
	hub := ws.NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWebSocket(w, r)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to dial WebSocket server: %v", err)
	}
	defer conn.Close()

	// Wait for broadcast metrics frame
	conn.SetReadDeadline(time.Now().Add(4 * time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("expected WebSocket frame from hub, got error: %v", err)
	}

	if !strings.Contains(string(msg), `"topic":"metrics"`) {
		t.Errorf("expected topic 'metrics' in message, got: %s", string(msg))
	}
}
