package service

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuditLog struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Action     string        `bson:"action" json:"action"`
	Method     string        `bson:"method" json:"method"`
	Path       string        `bson:"path" json:"path"`
	StatusCode int           `bson:"status_code" json:"status_code"`
	UserID     string        `bson:"user_id,omitempty" json:"user_id,omitempty"`
	UserEmail  string        `bson:"user_email,omitempty" json:"user_email,omitempty"`
	Role       string        `bson:"role,omitempty" json:"role,omitempty"`
	IPAddress  string        `bson:"ip_address" json:"ip_address"`
	UserAgent  string        `bson:"user_agent" json:"user_agent"`
	Timestamp  time.Time     `bson:"timestamp" json:"timestamp"`
	Details    string        `bson:"details,omitempty" json:"details,omitempty"`
}

type AuditLogService interface {
	Log(action, method, path string, statusCode int, userID, userEmail, role, ip, userAgent, details string) error
}

type auditLogService struct {
	col *mongo.Collection
}

func NewAuditLogService(db *mongo.Database) AuditLogService {
	return &auditLogService{
		col: db.Collection("audit_logs"),
	}
}

func (s *auditLogService) Log(action, method, path string, statusCode int, userID, userEmail, role, ip, userAgent, details string) error {
	audit := AuditLog{
		ID:         bson.NewObjectID(),
		Action:     action,
		Method:     method,
		Path:       path,
		StatusCode: statusCode,
		UserID:     userID,
		UserEmail:  userEmail,
		Role:       role,
		IPAddress:  ip,
		UserAgent:  userAgent,
		Timestamp:  time.Now(),
		Details:    details,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.col.InsertOne(ctx, audit)
	if err != nil {
		log.Printf("Error inserting audit log: %v", err)
	}
	return err
}
