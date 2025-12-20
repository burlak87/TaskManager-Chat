# Реализация WebSocket модуля для чата

## Выполненные задачи

### ✅ 1. WebSocket-сервер через gorilla/websocket
- Реализован WebSocket сервер с использованием библиотеки `gorilla/websocket`
- Поддержка множественных соединений через Hub паттерн
- Изоляция соединений по комнатам (доскам)

### ✅ 2. Обработка сообщений чата и broadcast
- Реализована обработка входящих сообщений от клиентов
- Автоматический broadcast сообщений всем клиентам в комнате
- Сохранение сообщений в MongoDB перед отправкой

### ✅ 3. Сохранение сообщений в БД MongoDB
- Создана модель `Message` для MongoDB
- Реализован репозиторий `MessageRepository` с методами:
  - `Create` - создание сообщения
  - `GetByBoardID` - получение сообщений для доски с пагинацией
  - `GetByID` - получение сообщения по ID
  - `CountByBoardID` - подсчет сообщений
  - `DeleteByBoardID` - удаление сообщений доски
- Интеграция с MongoDB через официальный драйвер

### ✅ 4. Управление комнатами (чат на доску)
- Каждая доска имеет свою комнату (room)
- Клиенты автоматически присоединяются к комнате при подключении
- Поддержка переключения между комнатами
- Изоляция сообщений между комнатами

### ✅ 5. Heartbeat, переподключение, обработка разрывов
- **Heartbeat**: Автоматическая отправка ping каждые 54 секунды
- **Pong handler**: Автоматический ответ на ping с обновлением deadline
- **Таймауты**: 
  - `pongWait` = 60 секунд (максимальное время ожидания ответа)
  - `pingPeriod` = 54 секунды (период отправки ping)
  - `writeWait` = 10 секунд (максимальное время записи)
- **Обработка разрывов**: Автоматическая очистка соединений при ошибках
- **Graceful shutdown**: Корректное закрытие соединений

## Структура модуля

```
backend/internal/websocket/
├── hub.go                 # Центральный hub для управления соединениями
├── client.go              # Клиент WebSocket с обработкой сообщений
├── message_processor.go   # Обработчик сообщений от клиентов
├── README.md              # Документация модуля
├── INTEGRATION.md         # Инструкция по интеграции
└── IMPLEMENTATION.md      # Этот файл

backend/internal/models/
├── message.go             # Модели для сообщений чата
└── websocket.go           # Модели для WebSocket сообщений

backend/internal/repository/
└── message_repository.go  # Репозиторий для работы с MongoDB

backend/internal/service/
└── chat_service.go        # Сервис для бизнес-логики чата

backend/internal/handler/
├── websocket_handler.go   # HTTP handler для WebSocket
└── chat_handler.go        # REST API handlers для чата

backend/pkg/database/
└── mongodb.go             # Подключение к MongoDB

backend/pkg/config/
└── config.go              # Конфигурация (включая MongoDB и WebSocket)
```

## Основные компоненты

### Hub
- Управляет всеми WebSocket соединениями
- Организует клиентов по комнатам (board_id)
- Обеспечивает thread-safe доступ к данным
- Обрабатывает регистрацию/отмену регистрации клиентов
- Выполняет broadcast сообщений в комнаты

### Client
- Представляет отдельное WebSocket соединение
- Обрабатывает входящие сообщения (readPump)
- Отправляет исходящие сообщения (writePump)
- Реализует heartbeat механизм
- Обрабатывает различные типы сообщений (message, ping, join, leave)

### MessageProcessor
- Обрабатывает сообщения от клиентов через канал
- Сохраняет сообщения в MongoDB
- Отправляет сообщения через WebSocket после сохранения
- Поддерживает callback для дополнительной обработки

## API Endpoints

### WebSocket
- `GET /ws/chat?board_id=<board_id>` - WebSocket подключение

### REST API
- `GET /api/boards/:board_id/messages` - Получить сообщения доски
- `GET /api/boards/:board_id/messages/count` - Количество сообщений
- `GET /api/messages/:message_id` - Получить сообщение по ID
- `POST /api/messages` - Создать сообщение

## Конфигурация

Переменные окружения:
- `MONGO_URI` - URI MongoDB (по умолчанию: `localhost:27017`)
- `MONGO_DATABASE` - Имя базы данных (по умолчанию: `taskmanager`)
- `MONGO_USER` - Пользователь MongoDB (опционально)
- `MONGO_PASSWORD` - Пароль MongoDB (опционально)

WebSocket настройки (в config):
- `ReadBufferSize`: 1024 байт
- `WriteBufferSize`: 1024 байт
- `PingPeriod`: 54 секунды
- `PongWait`: 60 секунд
- `WriteWait`: 10 секунд
- `MaxMessageSize`: 512 KB
- `EnableCompression`: true

## Безопасность

- Проверка аутентификации через middleware (ожидает `user_id` и `username` в контексте)
- Изоляция комнат (клиенты видят только сообщения своей доски)
- Валидация входящих сообщений
- Обработка ошибок с отправкой клиенту

## Производительность

- Буферизованные каналы для предотвращения блокировок
- Thread-safe операции через mutex
- Эффективное управление памятью (автоматическая очистка пустых комнат)
- Поддержка компрессии WebSocket

## Тестирование

Для тестирования можно использовать:
1. Браузерную консоль (см. INTEGRATION.md)
2. `websocat` утилиту
3. Postman или другие WebSocket клиенты

## Зависимости

- `github.com/gorilla/websocket` - WebSocket библиотека
- `go.mongodb.org/mongo-driver` - MongoDB драйвер
- `github.com/gin-gonic/gin` - HTTP фреймворк

## Следующие шаги

1. Интеграция с основным приложением (см. `examples/websocket_integration.go`)
2. Добавление middleware для JWT аутентификации
3. Тестирование в реальных условиях
4. Мониторинг и логирование (опционально)


