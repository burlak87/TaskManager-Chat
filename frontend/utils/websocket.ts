class WebSocketService {
  private socket: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectTimeout = 3000
  private messageHandlers: ((data: any) => void)[] = []
  private boardId: number = 0

  connect(boardId: number, token: string) {
    this.boardId = boardId
    const wsUrl = `ws://localhost:8888/api/ws/chat?board_id=${boardId}`
    
    this.socket = new WebSocket(wsUrl)
    
    this.socket.onopen = () => {
      console.log('WebSocket connected')
      this.reconnectAttempts = 0
    }
    
    this.socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        this.messageHandlers.forEach(handler => handler(data))
      } catch (error) {
        console.error('Error parsing message:', error)
      }
    }
    
    this.socket.onclose = (event) => {
      console.log('WebSocket disconnected:', event.code, event.reason)
      this.attemptReconnect()
    }
    
    this.socket.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
  }

  sendMessage(content: string) {
    if (this.socket?.readyState === WebSocket.OPEN) {
      const message = {
        board_id: this.boardId,
        content
      }
      this.socket.send(JSON.stringify(message))
    }
  }

  onMessage(handler: (data: any) => void) {
    this.messageHandlers.push(handler)
  }

  disconnect() {
    if (this.socket) {
      this.socket.close()
      this.socket = null
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`Reconnecting attempt ${this.reconnectAttempts}...`)
      
      setTimeout(() => {

        const token = localStorage.getItem('token') || ''
        this.connect(this.boardId, token)
      }, this.reconnectTimeout)
    }
  }
}

export const wsService = new WebSocketService()