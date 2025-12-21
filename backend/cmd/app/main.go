package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/your-team/taskmanager-chat/backend/internal/adapters/rest"
	"github.com/your-team/taskmanager-chat/backend/internal/service"
	mongodbstorage "github.com/your-team/taskmanager-chat/backend/internal/storage/mongodb"
	"github.com/your-team/taskmanager-chat/backend/internal/storage/psql"
	"github.com/your-team/taskmanager-chat/backend/internal/websocket"
	mongodbclient "github.com/your-team/taskmanager-chat/backend/pkg/client-database/mongodb"
	postgresqlclient "github.com/your-team/taskmanager-chat/backend/pkg/client-database/postgresql"
	"github.com/your-team/taskmanager-chat/backend/pkg/client/postgresql"
	"github.com/your-team/taskmanager-chat/backend/pkg/config"
	"github.com/your-team/taskmanager-chat/backend/pkg/logging"
	"github.com/your-team/taskmanager-chat/backend/pkg/middleware"
	"github.com/your-team/taskmanager-chat/backend/pkg/server"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Infoln("Starting application")
	
	cfg := config.GetConfig()
	logger.Infof("Environment: %s", cfg.Env)
	logger.Infof("DB CONFIG: Host=%s, Port=%s, Database=%s, Username=%s", cfg.Host, cfg.Port, cfg.Database, cfg.Username)
	
	postgresqSQLClient, err := postgresqlclient.NewClient(context.TODO(), 15, cfg.StorageConfig)
	if err != nil {
		logger.Fatlg("Failed to connect to database: %v", err)
	}
	defer postgresqSQLClient.Close()
	
	queries := sqlc.New(pgPool)
	storage := psql.NewStorage(queries)
	
	mongoClient, err := mongodbclient.NewClient(context.TODO(), 15, cfg.MongoConfig)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.TODO())
	
	messageStorage := mongodbstorage.NewMessageStorage(mongoClient, cfg.MongoConfig.Database)
	
	jwtSecret := "my-secret-key"
	logger.Infof("secret %s", jwtSecret)	
	
	userService := service.NewUser(storage, storage, jwtSecret)

	userHandler := rest.NewUsersHandler(userService, logger)
	
	wsHub := websocket.NewHub(messageStorage, logger.Logger)
	go wsHub.Run()
	
	wsHandler := websocket.NewHandler(wsHub, logger.Logger)

	serverCfg := server.Config{
		Port:         "8888",
		Mode:         cfg.Env,
		CorsOrigins:  []string{"*"},
		CorsEnabled:  true,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	
	srv := server.NewServer(serverCfg, logger.Logger)
	
	srv.RegisterRoutes(func(engine *gin.Engine) {
		api := engine.Group("/api")
		{
			auth := api.Group("/auth")
			userHandler.RegisterRoutes(auth, jwtSecret)
			
			ws := api.Group("/ws")
			ws.Use(middleware.JWTAuthMiddleware(jwtSecret))
			{
				ws.GET("/chat", wsHandler.HandleWebSocket)
			}
		}
		
		engine.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"time":   time.Now().UTC(),
			})
		})
	})
	
	if err := srv.Start(); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

func getDSN(cfg *config.Config) string {
	return "postgresql://" + cfg.Username + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.Database + "?sslmode=disable&pool_max_conns=20"
}
	
// postgreSQLClient, err := postgresql.NewClient(context.TODO(), 15, cfg.StorageConfig)
// 	if err != nil {
// 		logger.Fatalf("Failed to connect to database: %v", err)
// 	}
	
// 	logger.Infoln("Checking available databases...")
// 	rows, err := postgreSQLClient.Query(context.Background(), "SELECT datname FROM pg_database WHERE datistemplate = false;")
// 	if err == nil {
// 		defer rows.Close()
// 		for rows.Next() {
// 			var dbName string
// 			rows.Scan(&dbName)
// 			logger.Infof("Available database: %s", dbName)
// 		}
// 	}
	
// 	logger.Infoln("Checking tables in current database...")
// 	rows, err = postgreSQLClient.Query(context.Background(),
// 		"SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';")
// 	if err == nil {
// 		defer rows.Close()
// 		for rows.Next() {
// 			var tableName string
// 			rows.Scan(&tableName)
// 			logger.Infof("Table: %s", tableName)
// 		}
// 	}
	
// 	var actualUserCount int
// 	err = postgreSQLClient.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&actualUserCount)
// 	if err != nil {
// 		logger.Errorf("Count users failed: %v", err)
// 	} else {
// 		logger.Infof("Actual users in connected DB: %d", actualUserCount)

// 		rows, err := postgreSQLClient.Query(context.Background(), "SELECT id, email FROM users LIMIT 5")
// 		if err == nil {
// 			defer rows.Close()
// 			for rows.Next() {
// 				var id int64
// 				var email string
// 				rows.Scan(&id, &email)
// 				logger.Infof("   User: ID=%d, Email=%s", id, email)
// 			}
// 		}