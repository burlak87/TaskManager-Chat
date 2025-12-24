<template>
  <article :class="['message-item', { 'own-message': isOwnMessage }]">
    <header class="message-header">
      <span class="username">{{ message.username }}</span>
      <time class="timestamp" :datetime="message.created_at">
        {{ formatTime(message.created_at) }}
      </time>
    </header>
    <p class="message-content">
      <template v-for="(part, index) in parsedContent" :key="index">
        <span v-if="part.type === 'text'">{{ part.text }}</span>
        <UButton
          v-else
          size="xs"
          color="primary"
          variant="ghost"
          class="mention-button"
          @click="$emit('mention', part.username)"
        >
          @{{ part.username }}
        </UButton>
      </template>
    </p>
  </article>
</template>

<script setup lang="ts">
import type { Message } from '~/types/chat'

const props = defineProps<{
  message: Message
  currentUserId: number
}>()

const emit = defineEmits<{
  mention: [username: string]
}>()

const isOwnMessage = computed(() => props.message.user_id === props.currentUserId)

const parsedContent = computed(() => {
  const content = props.message.content
  const regex = /@(\w+)/g
  const parts = []
  let lastIndex = 0
  let match

  while ((match = regex.exec(content)) !== null) {
    if (match.index > lastIndex) {
      parts.push({
        type: 'text',
        text: content.substring(lastIndex, match.index)
      })
    }
    
    parts.push({
      type: 'mention',
      username: match[1]
    })
    
    lastIndex = match.index + match[0].length
  }

  if (lastIndex < content.length) {
    parts.push({
      type: 'text',
      text: content.substring(lastIndex)
    })
  }

  return parts
})

const formatTime = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.message-item {
  padding: 8px 12px;
  margin: 4px 0;
  border-radius: 8px;
  background: #f5f5f5;
  max-width: 80%;
}

.own-message {
  background: #e3f2fd;
  margin-left: auto;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
  font-size: 12px;
}

.username {
  font-weight: bold;
}

.timestamp {
  color: #666;
  font-size: 11px;
}

.message-content {
  word-break: break-word;
  margin: 0;
}

.mention-button {
  padding: 2px 4px;
  margin: 0 2px;
  vertical-align: baseline;
}
</style>