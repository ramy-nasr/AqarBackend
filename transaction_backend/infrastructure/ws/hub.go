package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"transaction-backend/domain"

	"github.com/gorilla/websocket"
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.Transaction, error)
}

type WebSocketHub struct {
	repo    Repository
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewWebSocketHub(repo Repository) *WebSocketHub {
	return &WebSocketHub{
		repo:    repo,
		clients: make(map[*websocket.Conn]bool),
	}
}

func (hub *WebSocketHub) HandleConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}

	// Send all past records first
	ctx := context.Background()
	records, err := hub.repo.GetAll(ctx)
	if err == nil {
		for _, txn := range records {
			msg, _ := json.Marshal(txn)
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}

	hub.mu.Lock()
	hub.clients[conn] = true
	hub.mu.Unlock()
}

func (hub *WebSocketHub) Broadcast(txn domain.Transaction) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	msg, _ := json.Marshal(txn)
	for conn := range hub.clients {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			conn.Close()
			delete(hub.clients, conn)
		}
	}
}
