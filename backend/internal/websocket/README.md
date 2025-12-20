# WebSocket модуль для чата

Этот модуль реализует WebSocket сервер для чата досок с использованием gorilla/websocket.

## Основные компоненты

### Hub
Центральный компонент для управления WebSocket соединениями и комнатами (досками).

### Client
Представляет отдельное WebSocket соединение с поддержкой:
- Heartbeat (ping/pong)
- Автоматическое переподключение
- Обработка разрывов соединения

### MessageProcessor
Обрабатывает сообщения от клиентов, сохраняет их в MongoDB и отправляет через WebSocket.

## Использование

### Инициализация

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

// Создаем hub
hub := websocket.NewHub()
go hub.Run()

// Создаем процессор сообщений
processor := websocket.NewMessageProcessor(hub, messageRepo)
processor.SetSaveCallback(func(msg *models.Message) (*models.MessageResponse, error) {
    // Здесь можно добавить дополнительную логику перед сохранением
    return chatService.SaveMessage(ctx, msg)
})
go processor.Start(ctx)

// Создаем handler
wsHandler := handler.NewWebSocketHandler(hub, cfg.WebSocket)
chatHandler := handler.NewChatHandler(chatService)

// Настраиваем роуты
router.GET("/ws/chat", wsHandler.HandleWebSocket)
router.GET("/api/boards/:board_id/messages", chatHandler.GetMessages)
router.POST("/api/messages", chatHandler.CreateMessage)
```

### Подключение клиента

Клиент подключается к `/ws/chat?board_id=<board_id>` с JWT токеном в заголовке Authorization.

### Формат сообщений

#### Входящие сообщения (от клиента)

```json
{
  "type": "message",
  "payload": {
    "content": "Текст сообщения",
    "board_id": "board_123",
    "mentions": ["user_1", "user_2"]
  }
}
```

#### Исходящие сообщения (к клиенту)

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

### Типы сообщений

- `message` - сообщение чата
- `ping` - ping для heartbeat
- `pong` - ответ на ping
- `join` - присоединение к комнате
- `leave` - выход из комнаты
- `error` - сообщение об ошибке

## Особенности

- **Heartbeat**: Автоматическая отправка ping каждые 54 секунды
- **Переподключение**: Клиент может переподключиться и продолжить работу
- **Комнаты**: Каждая доска имеет свою комнату для изоляции сообщений
- **Broadcast**: Сообщения автоматически отправляются всем клиентам в комнате
- **MongoDB**: Сообщения сохраняются в MongoDB для истории


