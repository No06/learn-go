package main

import (
	"fmt"

	"hinoob.net/learn-go/internal/config"
	"hinoob.net/learn-go/internal/pkg/websocket"
	"hinoob.net/learn-go/internal/repository"
	"hinoob.net/learn-go/internal/router"
	"hinoob.net/learn-go/pkg/oss"
)

func main() {
	// Load configuration
	config.LoadConfig("./configs")

	// Initialize database
	repository.InitDB()

	// Initialize OSS Client
	oss.InitOSS()

	// Create and run the websocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Set up router
	r := router.SetupRouter(hub)

	// Start the server
	port := config.AppConfig.Server.Port
	fmt.Printf("Server is running on port %s\n", port)
	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
