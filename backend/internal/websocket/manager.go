package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/smartstocks/backend/internal/models"
)

type Client struct {
	ID      string
	UserID  string
	Conn    *websocket.Conn
	Send    chan []byte
	Manager *Manager
	MatchID string
	mu      sync.Mutex
}

type Manager struct {
	clients    map[string]*Client
	matches    map[string][]*Client
	Register   chan *Client // Exportado
	Unregister chan *Client // Exportado
	broadcast  chan *BroadcastMessage
	mu         sync.RWMutex
}

type BroadcastMessage struct {
	MatchID string
	Message []byte
	Exclude string
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[string]*Client),
		matches:    make(map[string][]*Client),
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.Register:
			m.registerClient(client)
		case client := <-m.Unregister:
			m.unregisterClient(client)
		case message := <-m.broadcast:
			m.broadcastToMatch(message)
		}
	}
}

func (m *Manager) registerClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[client.UserID] = client
	log.Printf("✅ Client registered: UserID=%s, Total clients=%d", client.UserID, len(m.clients))
}

func (m *Manager) unregisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.clients[client.UserID]; ok {
		delete(m.clients, client.UserID)
		close(client.Send)

		if client.MatchID != "" {
			m.removeFromMatch(client)
		}

		log.Printf("❌ Client unregistered: UserID=%s, Total clients=%d", client.UserID, len(m.clients))
	}
}

func (m *Manager) AddToMatch(client *Client, matchID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	client.MatchID = matchID
	m.matches[matchID] = append(m.matches[matchID], client)

	log.Printf("✅ Client added to match: UserID=%s, MatchID=%s, Players in match=%d",
		client.UserID, matchID, len(m.matches[matchID]))
}

func (m *Manager) removeFromMatch(client *Client) {
	if clients, ok := m.matches[client.MatchID]; ok {
		for i, c := range clients {
			if c.UserID == client.UserID {
				m.matches[client.MatchID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}

		if len(m.matches[client.MatchID]) == 0 {
			delete(m.matches, client.MatchID)
		}
	}
}

func (m *Manager) BroadcastToMatch(matchID string, message interface{}, excludeUserID string) error {
	wsMessage := models.WSMessage{
		Type:      getMessageType(message),
		Data:      message,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		return err
	}

	m.broadcast <- &BroadcastMessage{
		MatchID: matchID,
		Message: data,
		Exclude: excludeUserID,
	}

	return nil
}

func (m *Manager) broadcastToMatch(msg *BroadcastMessage) {
	m.mu.RLock()
	clients := m.matches[msg.MatchID]
	m.mu.RUnlock()

	log.Printf("📤 Broadcasting to %d clients in match %s", len(clients), msg.MatchID)

	for _, client := range clients {
		if msg.Exclude != "" && client.UserID == msg.Exclude {
			continue
		}

		select {
		case client.Send <- msg.Message:
			log.Printf("✅ Message sent to user %s", client.UserID)
		default:
			log.Printf("❌ Failed to send to user %s, closing connection", client.UserID)
			close(client.Send)
			m.Unregister <- client
		}
	}
}

func (m *Manager) SendToUser(userID string, message interface{}) error {
	m.mu.RLock()
	client, exists := m.clients[userID]
	m.mu.RUnlock()

	if !exists {
		log.Printf("⚠️  User %s not connected, cannot send message", userID)
		return nil
	}

	wsMessage := models.WSMessage{
		Type:      getMessageType(message),
		Data:      message,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(wsMessage)
	if err != nil {
		return err
	}

	select {
	case client.Send <- data:
		log.Printf("✅ Message sent to user %s", userID)
	default:
		log.Printf("❌ Failed to send to user %s", userID)
		close(client.Send)
		m.Unregister <- client
	}

	return nil
}

func (m *Manager) GetClient(userID string) (*Client, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, exists := m.clients[userID]
	return client, exists
}

func (m *Manager) IsUserConnected(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.clients[userID]
	return exists
}

func (m *Manager) GetMatchClients(matchID string) []*Client {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.matches[matchID]
}

// === CLIENT METHODS ===

func (c *Client) ReadPump() {
	defer func() {
		c.Manager.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		c.handleMessage(message)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(message []byte) {
	var wsMsg models.WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	switch wsMsg.Type {
	case models.WSMsgTypePing:
		c.SendMessage(models.WSMsgTypePong, nil)
	default:
		log.Printf("Unknown message type: %s", wsMsg.Type)
	}
}

func (c *Client) SendMessage(msgType models.WSMessageType, data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	wsMessage := models.WSMessage{
		Type:      msgType,
		Data:      data,
		Timestamp: time.Now(),
	}

	messageData, err := json.Marshal(wsMessage)
	if err != nil {
		return err
	}

	select {
	case c.Send <- messageData:
	default:
		return nil
	}

	return nil
}

func (c *Client) SendError(errorMsg string) error {
	return c.SendMessage(models.WSMsgTypeError, map[string]string{
		"error": errorMsg,
	})
}

func getMessageType(message interface{}) models.WSMessageType {
	switch message.(type) {
	case *models.MatchFoundResponse:
		return models.WSMsgTypeMatchFound
	case *models.RoundStartResponse:
		return models.WSMsgTypeRoundStart
	case *models.RoundResultResponse:
		return models.WSMsgTypeRoundResult
	case *models.MatchResultResponse:
		return models.WSMsgTypeMatchResult
	default:
		return models.WSMsgTypeError
	}
}
