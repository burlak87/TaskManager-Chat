<template>
  <UApp>
    <header class="app-header">
      <nav class="app-nav">
        <h1>TaskManager Chat</h1>
        <section class="user-info">
          <UBadge color="blue" variant="soft">
            User ID: {{ currentUserId }}
          </UBadge>
          <UButton
            v-if="!isLoggedIn"
            color="primary"
            @click="login"
          >
            Login (Demo)
          </UButton>
        </section>
      </nav>
    </header>
    
    <main class="app-main">
      <ChatWindow :boardId="1" />
    </main>
  </UApp>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ChatWindow from '~/components/ChatWindow.vue'

const currentUserId = ref(0)
const isLoggedIn = ref(false)

const login = () => {
  currentUserId.value = Math.floor(Math.random() * 1000) + 1
  localStorage.setItem('userId', currentUserId.value.toString())
  localStorage.setItem('token', 'demo-token-' + Date.now())
  localStorage.setItem('username', 'User' + currentUserId.value)
  isLoggedIn.value = true
}

onMounted(() => {
  const userId = localStorage.getItem('userId')
  if (userId) {
    currentUserId.value = parseInt(userId)
    isLoggedIn.value = true
  }
})
</script>

<style scoped>
.app-header {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  background: white;
}

.app-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
}

.app-nav h1 {
  font-size: 1.5rem;
  font-weight: bold;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-main {
  max-width: 1200px;
  margin: 0 auto;
  height: calc(100vh - 73px);
}
</style>