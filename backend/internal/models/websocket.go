package models

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type IncomingMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type OutgoingMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type JoinRoomPayload struct {
	BoardID string `json:"board_id"`
}

type LeaveRoomPayload struct {
	BoardID string `json:"board_id"`
}

type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type PingPayload struct {
	Timestamp int64 `json:"timestamp"`
}

type PongPayload struct {
	Timestamp int64 `json:"timestamp"`
}

