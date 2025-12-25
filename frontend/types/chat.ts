export interface Message {
  id: string
  board_id: number
  user_id: number
  username: string
  content: string
  created_at: string
}

export interface MessageRequest {
  board_id: number
  content: string
}

export interface Participant {
  id: number
  username: string
  isOnline: boolean
}