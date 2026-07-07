package service

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"go-tiket-konser/models"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"gorm.io/gorm"
)

type ChatMessage struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID     string        `bson:"room_id" json:"room_id"` // Matches Customer User ID
	SenderID   string        `bson:"sender_id" json:"sender_id"`
	SenderName string        `bson:"sender_name" json:"sender_name"`
	Role       string        `bson:"role" json:"role"` // customer or admin
	Message    string        `bson:"message" json:"message"`
	Timestamp  time.Time     `bson:"timestamp" json:"timestamp"`
}

type ChatRoomResponse struct {
	RoomID       string      `json:"room_id"`
	CustomerName string      `json:"customer_name"`
	LastMessage  ChatMessage `json:"last_message"`
}

type Client struct {
	UserID     string
	FullName   string
	Role       string
	RoomID     string
	Conn       *websocket.Conn
	Send       chan []byte
	Hub        *ChatHub
}

type ChatHub struct {
	mu       sync.RWMutex
	rooms    map[string]map[*Client]bool // roomID -> clients set
	col      *mongo.Collection
	postgres *gorm.DB
}

func NewChatHub(mongoDB *mongo.Database, pgDB *gorm.DB) *ChatHub {
	return &ChatHub{
		rooms:    make(map[string]map[*Client]bool),
		col:      mongoDB.Collection("chat_messages"),
		postgres: pgDB,
	}
}

func (h *ChatHub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[c.RoomID] == nil {
		h.rooms[c.RoomID] = make(map[*Client]bool)
	}
	h.rooms[c.RoomID][c] = true
	log.Printf("Client registered to room %s: %s (%s)", c.RoomID, c.FullName, c.Role)
}

func (h *ChatHub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[c.RoomID] != nil {
		if _, ok := h.rooms[c.RoomID][c]; ok {
			delete(h.rooms[c.RoomID], c)
			close(c.Send)
			if len(h.rooms[c.RoomID]) == 0 {
				delete(h.rooms, c.RoomID)
			}
		}
	}
	log.Printf("Client unregistered from room %s: %s (%s)", c.RoomID, c.FullName, c.Role)
}

func (h *ChatHub) Broadcast(roomID string, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if clients, exists := h.rooms[roomID]; exists {
		for c := range clients {
			select {
			case c.Send <- msg:
			default:
				// If client channel is blocked, handle unregister later in write pump
			}
		}
	}
}

func (h *ChatHub) SaveMessage(msg ChatMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := h.col.InsertOne(ctx, msg)
	return err
}

func (h *ChatHub) GetRoomMessages(roomID string) ([]ChatMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"room_id": roomID}
	opts := options.Find().SetSort(bson.M{"timestamp": 1})
	cursor, err := h.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []ChatMessage = make([]ChatMessage, 0)
	for cursor.Next(ctx) {
		var msg ChatMessage
		if err := cursor.Decode(&msg); err != nil {
			continue
		}
		list = append(list, msg)
	}
	return list, nil
}

func (h *ChatHub) GetRoomsHistory() ([]ChatRoomResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Aggregation pipeline to get unique rooms and their last message
	pipeline := []bson.M{
		{"$sort": bson.M{"timestamp": -1}},
		{"$group": bson.M{
			"_id":          "$room_id",
			"last_message": bson.M{"$first": "$$ROOT"},
		}},
	}

	cursor, err := h.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []ChatRoomResponse = make([]ChatRoomResponse, 0)
	for cursor.Next(ctx) {
		var result struct {
			RoomID      string      `bson:"_id"`
			LastMessage ChatMessage `bson:"last_message"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}

		// Retrieve Customer Name from Postgres (Room ID matches Customer User ID)
		customerName := "Unknown Customer"
		var user models.User
		if err := h.postgres.Preload("Customer").First(&user, "id = ?", result.RoomID).Error; err == nil {
			if user.Customer != nil {
				customerName = user.Customer.Name
			} else {
				customerName = user.FullName
			}
		}

		list = append(list, ChatRoomResponse{
			RoomID:       result.RoomID,
			CustomerName: customerName,
			LastMessage:  result.LastMessage,
		})
	}
	return list, nil
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(4096) // Set max message size
	_ = c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading WebSocket: %v", err)
			}
			break
		}

		// Parse received JSON message
		var input struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(message, &input); err != nil {
			continue
		}

		if input.Message == "" {
			continue
		}

		chatMsg := ChatMessage{
			ID:         bson.NewObjectID(),
			RoomID:     c.RoomID,
			SenderID:   c.UserID,
			SenderName: c.FullName,
			Role:       c.Role,
			Message:    input.Message,
			Timestamp:  time.Now(),
		}

		// Save message to MongoDB
		if err := c.Hub.SaveMessage(chatMsg); err != nil {
			log.Printf("error saving chat message: %v", err)
		}

		// Broadcast message to room
		payload, err := json.Marshal(chatMsg)
		if err == nil {
			c.Hub.Broadcast(c.RoomID, payload)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
