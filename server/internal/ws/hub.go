package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: restrict to allowed dashboard origin in production
		return true
	},
}

// Hub manages all active WebSocket connections and broadcasts metric updates.
type Hub struct {
	mu         sync.RWMutex
	clients    map[*client]struct{}
	broadcast  chan []byte
	register   chan *client
	unregister chan *client
}

type client struct {
	conn *websocket.Conn
	send chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*client]struct{}),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *client, 16),
		unregister: make(chan *client, 16),
	}
}

// Run is the hub's main event loop. Call as a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					// Slow client — drop message to avoid blocking
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a JSON payload to all connected dashboard clients.
func (h *Hub) Broadcast(payload []byte) {
	select {
	case h.broadcast <- payload:
	default:
		log.Println("ws hub: broadcast channel full, dropping message")
	}
}

// ServeWS is the Gin handler that upgrades HTTP to WebSocket.
func (h *Hub) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	cl := &client{conn: conn, send: make(chan []byte, 64)}
	h.register <- cl

	go cl.writePump(h)
	go cl.readPump(h)
}

func (c *client) writePump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (c *client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()
	for {
		// Dashboard only reads (no commands from client yet); drain to detect disconnects
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}
