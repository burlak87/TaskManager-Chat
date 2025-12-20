package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	WebSocket WebSocketConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type MongoDBConfig struct {
	URI      string
	Database string
	Username string
	Password string
}

type WebSocketConfig struct {
	ReadBufferSize    int
	WriteBufferSize   int
	PingPeriod        int
	PongWait          int
	WriteWait         int
	MaxMessageSize    int64
	EnableCompression bool
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGO_URI", "localhost:27017"),
			Database: getEnv("MONGO_DATABASE", "taskmanager"),
			Username: getEnv("MONGO_USER", ""),
			Password: getEnv("MONGO_PASSWORD", ""),
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			PingPeriod:        54,
			PongWait:          60,
			WriteWait:         10,
			MaxMessageSize:    512 * 1024,
			EnableCompression: true,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

