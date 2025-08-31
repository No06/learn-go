package handler

import (
	"hinoob.net/learn-go/internal/pkg/jwt"
	"log"
	"net/http"

	"hinoob.net/learn-go/internal/pkg/websocket"

	"github.com/gin-gonic/gin"
	gorillaWebsocket "github.com/gorilla/websocket"
)

var wsUpgrader = gorillaWebsocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, check the origin to prevent CSRF attacks.
		// Example: return r.Header.Get("Origin") == "http://your-frontend.com"
		return true
	},
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *websocket.Hub, c *gin.Context) {
	// 1. Authenticate the user from the token in the query param
	tokenString := c.Query("token")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
		return
	}

	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// 2. Upgrade the HTTP connection to a WebSocket connection
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// 3. Create a new client and register it with the hub
	client := &websocket.Client{
		Hub:    hub,
		UserID: claims.UserID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	client.Hub.Register <- client

	// 4. Start the client's read and write pumps in separate goroutines
	go client.WritePump()
	go client.ReadPump()
}
