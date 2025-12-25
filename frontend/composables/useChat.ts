import { ref, computed, onUnmounted } from 'vue'
import type { Message, Participant } from '~/types/chat'
import { wsService } from '~/utils/websocket'
import { useToast } from '#imports'

export const useChat = (boardId: number) => {
  const messages = ref<Message[]>([])
  const participants = ref<Participant[]>([
    { id: 1, username: 'User1', isOnline: true },
    { id: 2, username: 'User2', isOnline: true },
    { id: 3, username: 'User3', isOnline: false }
  ])
  const inputMessage = ref('')
  const isLoading = ref(false)
  const toast = useToast()

  const connect = () => {
    const token = localStorage.getItem('token')
    if (!token) {
      toast.add({ title: 'Authentication required', color: 'red' })
      return
    }

    wsService.connect(boardId, token)
    
    wsService.onMessage((data) => {
      if (data.board_id === boardId) {
        messages.value.push(data)
        
        if (data.user_id !== parseInt(localStorage.getItem('userId') || '0')) {
          toast.add({
            title: `${data.username}`,
            description: data.content,
            timeout: 3000
          })
        }
        
        scrollToBottom()
      }
    })
  }

  const sendMessage = () => {
    if (!inputMessage.value.trim()) return
    
    wsService.sendMessage(inputMessage.value)
    inputMessage.value = ''
  }

  const scrollToBottom = () => {
    nextTick(() => {
      const container = document.querySelector('.messages-container')
      if (container) {
        container.scrollTop = container.scrollHeight
      }
    })
  }

  const formatTime = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  }

  const handleMention = (username: string) => {
    inputMessage.value = inputMessage.value + `@${username} `
  }

  onUnmounted(() => {
    wsService.disconnect()
  })

  return {
    messages,
    participants,
    inputMessage,
    isLoading,
    connect,
    sendMessage,
    formatTime,
    handleMention
  }
}