# Интеграция WebSocket модуля

## Быстрый старт

### 1. Установка зависимостей

```bash
cd backend
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
go get go.mongodb.org/mongo-driver/mongo/options
go mod tidy
```

### 2. Настройка переменных окружения

Добавьте в `.env` или переменные окружения:

```env
MONGO_URI=mongodb:27017
MONGO_DATABASE=taskmanager
MONGO_USER=admin
MONGO_PASSWORD=password
```

### 3. Инициализация в коде

```go
import (
    "context"
    "github.com/your-team/taskmanager-chat/backend/internal/websocket"
    "github.com/your-team/taskmanager-chat/backend/internal/repository"
    "github.com/your-team/taskmanager-chat/backend/internal/service"
    "github.com/your-team/taskmanager-chat/backend/pkg/database"
    "github.com/your-team/taskmanager-chat/backend/pkg/config"
)

// Загружаем конфигурацию
cfg, _ := config.Load()

// Подключаемся к MongoDB
mongoDB, _ := database.NewMongoDB(ctx, database.Config{
    URI:      cfg.MongoDB.URI,
    Database: cfg.MongoDB.Database,
    Username: cfg.MongoDB.Username,
    Password: cfg.MongoDB.Password,
})
defer mongoDB.Close(ctx)

// Создаем репозиторий
messageRepo := repository.NewMessageRepository(mongoDB.Database)

// Создаем hub и запускаем его
hub := websocket.NewHub()
go hub.Run()

// Создаем сервис
chatService := service.NewChatService(messageRepo, hub)

// Создаем процессор сообщений
processor := websocket.NewMessageProcessor(hub, messageRepo)
processor.SetSaveCallback(func(msg *models.Message) (*models.MessageResponse, error) {
    return chatService.SaveMessage(ctx, msg)
})
go processor.Start(ctx)

// Создаем handlers
wsHandler := handler.NewWebSocketHandler(hub, cfg.WebSocket)
chatHandler := handler.NewChatHandler(chatService)
```

### 4. Настройка роутов

```go
// WebSocket endpoint
router.GET("/ws/chat", authMiddleware, wsHandler.HandleWebSocket)

// REST API endpoints
api := router.Group("/api")
{
    api.GET("/boards/:board_id/messages", authMiddleware, chatHandler.GetMessages)
    api.POST("/messages", authMiddleware, chatHandler.CreateMessage)
}
```

## Middleware для аутентификации

WebSocket handler ожидает, что middleware установит в контекст:
- `user_id` (string) - ID пользователя
- `username` (string) - имя пользователя

Пример middleware:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        
        // Проверяем JWT токен и извлекаем данные пользователя
        claims, err := validateJWT(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

## Подключение клиента

### WebSocket URL

```
ws://localhost:8080/ws/chat?board_id=<board_id>
```

### Заголовки

```
Authorization: Bearer <jwt_token>
```

### Пример подключения (JavaScript)

```javascript
const token = 'your_jwt_token';
const boardId = 'board_123';

const ws = new WebSocket(`ws://localhost:8080/ws/chat?board_id=${boardId}`, {
    headers: {
        'Authorization': `Bearer ${token}`
    }
});

ws.onopen = () => {
    console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('Received:', message);
};

// Отправка сообщения
ws.send(JSON.stringify({
    type: 'message',
    payload: {
        content: 'Hello, world!',
        board_id: boardId,
        mentions: ['user_1', 'user_2']
    }
}));
```

## Формат сообщений

### Входящие (от клиента)

```json
{
  "type": "message",
  "payload": {
    "content": "Текст сообщения",
    "board_id": "board_123",
    "mentions": ["user_1"]
  }
}
```

### Исходящие (к клиенту)

```json
{
  "type": "message",
  "payload": {
    "id": "message_id",
    "board_id": "board_123",
    "user_id": "user_123",
    "username": "Имя пользователя",
    "content": "Текст сообщения",
    "mentions": ["user_1"],
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

## Типы сообщений

- `message` - сообщение чата
- `ping` - ping для heartbeat (клиент → сервер)
- `pong` - ответ на ping (сервер → клиент)
- `join` - присоединение к комнате
- `leave` - выход из комнаты
- `joined` - подтверждение присоединения
- `left` - подтверждение выхода
- `error` - сообщение об ошибке

## Особенности

1. **Heartbeat**: Автоматическая отправка ping каждые 54 секунды
2. **Автоматическое переподключение**: Клиент может переподключиться
3. **Комнаты**: Каждая доска имеет свою комнату
4. **Broadcast**: Сообщения автоматически отправляются всем в комнате
5. **Сохранение в БД**: Все сообщения сохраняются в MongoDB

## Тестирование

Для тестирования можно использовать `websocat`:

```bash
websocat "ws://localhost:8080/ws/chat?board_id=test_board" \
  -H "Authorization: Bearer your_token"
```

Или через браузерную консоль (см. пример выше).


