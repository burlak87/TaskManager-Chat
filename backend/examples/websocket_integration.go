package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/your-team/taskmanager-chat/backend/internal/handler"
	"github.com/your-team/taskmanager-chat/backend/internal/models"
	"github.com/your-team/taskmanager-chat/backend/internal/repository"
	"github.com/your-team/taskmanager-chat/backend/internal/service"
	"github.com/your-team/taskmanager-chat/backend/internal/websocket"
	"github.com/your-team/taskmanager-chat/backend/pkg/config"
	"github.com/your-team/taskmanager-chat/backend/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mongoDB, err := database.NewMongoDB(ctx, database.Config{
		URI:      cfg.MongoDB.URI,
		Database: cfg.MongoDB.Database,
		Username: cfg.MongoDB.Username,
		Password: cfg.MongoDB.Password,
	})

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close(ctx)

	log.Println("Connected to MongoDB successfully")
	messageRepo := repository.NewMessageRepository(mongoDB.Database)
	hub := websocket.NewHub()
	go hub.Run()
	log.Println("WebSocket hub started")
	chatService := service.NewChatService(messageRepo, hub)
	processor := websocket.NewMessageProcessor(hub, messageRepo)
	processor.SetSaveCallback(func(msg *models.Message) (*models.MessageResponse, error) {
		return chatService.SaveMessage(ctx, msg)
	})
	go processor.Start(ctx)
	log.Println("Message processor started")
	wsHandler := handler.NewWebSocketHandler(hub, cfg.WebSocket)
	chatHandler := handler.NewChatHandler(chatService)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/ws/chat", wsHandler.HandleWebSocket)

	api := router.Group("/api")
	{
		api.GET("/boards/:board_id/messages", chatHandler.GetMessages)
		api.GET("/boards/:board_id/messages/count", chatHandler.GetMessagesCount)
		api.GET("/messages/:message_id", chatHandler.GetMessageByID)
		api.POST("/messages", chatHandler.CreateMessage)
	}

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", serverAddr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := router.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down server...")
	cancel()
	log.Println("Server stopped")
}
