<template>
  <article class="chat-container">
    <!-- Participants Sidebar -->
    <aside class="sidebar">
      <header class="sidebar-header">
        <h3>Participants</h3>
        <UBadge color="green" variant="soft">
          {{ onlineCount }} online
        </UBadge>
      </header>
      
      <nav class="participants-list">
        <section
          v-for="participant in participants"
          :key="participant.id"
          class="participant-item"
        >
          <figure class="participant-info">
            <figcaption class="avatar">
              {{ participant.username.charAt(0) }}
            </figcaption>
            <span>{{ participant.username }}</span>
          </figure>
          <time :class="['status', participant.isOnline ? 'online' : 'offline']" :datetime="participant.isOnline ? 'online' : 'offline'"></time>
        </section>
      </nav>
    </aside>

    
    <main class="main-chat">
      <header class="chat-header">
        <h2>Chat Room #{{ boardId }}</h2>
        <UButton
          v-if="!isConnected"
          color="primary"
          @click="connect"
          :loading="isLoading"
        >
          Connect
        </UButton>
        <UBadge v-else color="green" variant="soft">
          Connected
        </UBadge>
      </header>

      
      <article ref="messagesContainer" class="messages-container">
        <section v-if="messages.length === 0" class="empty-state">
          <UIcon name="i-heroicons-chat-bubble-left-right" class="empty-icon" />
          <p>No messages yet. Start the conversation!</p>
        </section>
        
        <ChatMessage
          v-for="message in messages"
          :key="message.id"
          :message="message"
          :current-user-id="currentUserId"
          @mention="handleMention"
        />
      </article>

      
      <footer class="message-input">
        <UInput
          v-model="inputMessage"
          placeholder="Type your message..."
          size="lg"
          :disabled="!isConnected"
          @keyup.enter="sendMessage"
        >
          <template #trailing>
            <UButton
              color="primary"
              :disabled="!inputMessage.trim() || !isConnected"
              @click="sendMessage"
            >
              Send
            </UButton>
          </template>
        </UInput>
        
        <section class="quick-mentions">
          <span class="label">Quick mentions:</span>
          <UButton
            v-for="participant in onlineParticipants"
            :key="participant.id"
            size="xs"
            color="gray"
            variant="ghost"
            @click="handleMention(participant.username)"
          >
            @{{ participant.username }}
          </UButton>
        </section>
      </footer>
    </main>
  </article>
</template>

<script setup lang="ts">
import { onMounted, computed, ref, nextTick, watch } from 'vue'
import { useChat } from '~/composables/useChat'
import ChatMessage from './ChatMessage.vue'

const props = defineProps<{
  boardId: number
}>()

const currentUserId = ref(parseInt(localStorage.getItem('userId') || '1'))
const isConnected = ref(false)
const messagesContainer = ref<HTMLElement>()

const {
  messages,
  participants,
  inputMessage,
  isLoading,
  connect: connectChat,
  sendMessage: sendChatMessage,
  handleMention: mentionUser
} = useChat(props.boardId)

const onlineCount = computed(() => 
  participants.value.filter(p => p.isOnline).length
)

const onlineParticipants = computed(() =>
  participants.value.filter(p => p.isOnline && p.id !== currentUserId.value)
)

const connect = () => {
  connectChat()
  isConnected.value = true
}

const sendMessage = () => {
  sendChatMessage()
}

const handleMention = (username: string) => {
  mentionUser(username)
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

onMounted(() => {
  setTimeout(() => {
    scrollToBottom()
  }, 100)
})

watch(messages, () => {
  scrollToBottom()
}, { deep: true })
</script>

<style scoped>
.chat-container {
  display: flex;
  height: 100vh;
  background: white;
}

.sidebar {
  width: 250px;
  border-right: 1px solid #e5e7eb;
  padding: 16px;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.participants-list {
  flex: 1;
  overflow-y: auto;
}

.participant-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  border-radius: 6px;
  margin-bottom: 4px;
  transition: background-color 0.2s;
}

.participant-item:hover {
  background-color: #f9fafb;
}

.participant-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
}

.status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status.online {
  background-color: #10b981;
}

.status.offline {
  background-color: #9ca3af;
}

.main-chat {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chat-header {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.messages-container {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #9ca3af;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.message-input {
  padding: 16px;
  border-top: 1px solid #e5e7eb;
}

.quick-mentions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
}

.label {
  font-size: 12px;
  color: #6b7280;
}
</style>