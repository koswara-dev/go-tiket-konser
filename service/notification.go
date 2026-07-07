package service

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Notification struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string        `bson:"user_id" json:"user_id"` // empty means broadcast
	Title     string        `bson:"title" json:"title"`
	Message   string        `bson:"message" json:"message"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

type NotificationBroker struct {
	mu      sync.RWMutex
	clients map[string]map[chan string]bool // userID -> set of client channels
	db      *mongo.Database
	col     *mongo.Collection
}

func NewNotificationBroker(db *mongo.Database) *NotificationBroker {
	broker := &NotificationBroker{
		clients: make(map[string]map[chan string]bool),
		db:      db,
		col:     db.Collection("notifications"),
	}
	return broker
}

func (b *NotificationBroker) Register(userID string, ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.clients[userID] == nil {
		b.clients[userID] = make(map[chan string]bool)
	}
	b.clients[userID][ch] = true
}

func (b *NotificationBroker) Unregister(userID string, ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.clients[userID] != nil {
		delete(b.clients[userID], ch)
		if len(b.clients[userID]) == 0 {
			delete(b.clients, userID)
		}
	}
}

func (b *NotificationBroker) SendNotification(userID string, title, message string) error {
	notification := Notification{
		ID:        bson.NewObjectID(),
		UserID:    userID,
		Title:     title,
		Message:   message,
		CreatedAt: time.Now(),
	}

	// Save to Mongo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := b.col.InsertOne(ctx, notification)
	if err != nil {
		log.Printf("Error saving notification to Mongo: %v", err)
	}

	// Send to SSE clients
	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	payloadStr := string(payload)

	b.mu.RLock()
	defer b.mu.RUnlock()

	// If targeted notification, send to the user's active channels
	if userID != "" {
		if channels, exists := b.clients[userID]; exists {
			for ch := range channels {
				select {
				case ch <- payloadStr:
				default:
					// channel full or blocked
				}
			}
		}
	} else {
		// Broadcast to everyone
		for _, channels := range b.clients {
			for ch := range channels {
				select {
				case ch <- payloadStr:
				default:
				}
			}
		}
	}

	return nil
}

func (b *NotificationBroker) GetHistory(userID string) ([]Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch notifications for this user OR broadcast notifications
	filter := bson.M{
		"$or": []bson.M{
			{"user_id": userID},
			{"user_id": ""},
		},
	}

	cursor, err := b.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []Notification = make([]Notification, 0)
	for cursor.Next(ctx) {
		var n Notification
		if err := cursor.Decode(&n); err != nil {
			continue
		}
		list = append(list, n)
	}
	return list, nil
}
