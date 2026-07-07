package handler

import (
	"log"
	"net/http"
	"strings"

	"go-tiket-konser/dto"
	"go-tiket-konser/models"
	"go-tiket-konser/repository"
	"go-tiket-konser/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ChatHandler struct {
	hub           *service.ChatHub
	blacklistRepo repository.BlacklistedTokenRepository
	postgres      *gorm.DB
}

func NewChatHandler(
	hub *service.ChatHub,
	blacklistRepo repository.BlacklistedTokenRepository,
	pgDB *gorm.DB,
) *ChatHandler {
	return &ChatHandler{
		hub:           hub,
		blacklistRepo: blacklistRepo,
		postgres:      pgDB,
	}
}

var chatUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow cross-origin connection upgrade
	},
}

func (h *ChatHandler) WS(c *gin.Context) {
	// Extract token from query parameter or header
	tokenString := c.Query("token")
	if tokenString == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		}
	}

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Validate token blacklisting
	blacklisted, err := h.blacklistRepo.IsTokenBlacklisted(tokenString)
	if err != nil || blacklisted {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted or invalid"})
		return
	}

	// Parse claims
	claims := &service.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return service.JWTSecretKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		return
	}

	// Retrieve user display name from Postgres
	var user models.User
	fullName := "Customer"
	if err := h.postgres.First(&user, "id = ?", claims.UserID).Error; err == nil {
		fullName = user.FullName
	}

	// Determine chat room
	var roomID string
	if claims.Role == "admin" {
		roomID = c.Query("room_id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "room_id query parameter is required for admin role"})
			return
		}
	} else {
		// Customers can only join their own room (identified by their user ID)
		roomID = claims.UserID.String()
	}

	// Upgrade to WebSocket connection
	conn, err := chatUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade chat connection to WebSocket: %v", err)
		return
	}

	client := &service.Client{
		UserID:     claims.UserID.String(),
		FullName:   fullName,
		Role:       claims.Role,
		RoomID:     roomID,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		Hub:        h.hub,
	}

	h.hub.Register(client)

	// Start pump loops
	go client.WritePump()
	go client.ReadPump()
}

func (h *ChatHandler) GetRooms(c *gin.Context) {
	rooms, err := h.hub.GetRoomsHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil daftar room chat",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Daftar room chat berhasil diambil",
		Data:    rooms,
	})
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	roomID := c.Param("roomId")

	// IDOR Protection: Customers can only fetch their own messages
	tokenUserIDVal, _ := c.Get("user_id")
	tokenUserID := tokenUserIDVal.(uuid.UUID).String()
	tokenRole := c.MustGet("role").(string)

	if tokenRole != "admin" && tokenUserID != roomID {
		c.JSON(http.StatusForbidden, dto.WebResponse{
			Success: false,
			Message: "Akses ditolak: Anda tidak dapat mengakses percakapan ini",
		})
		return
	}

	messages, err := h.hub.GetRoomMessages(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.WebResponse{
			Success: false,
			Message: "Gagal mengambil riwayat pesan",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebResponse{
		Success: true,
		Message: "Riwayat pesan berhasil diambil",
		Data:    messages,
	})
}
