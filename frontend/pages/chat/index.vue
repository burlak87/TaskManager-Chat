<template>
  <div class="h-[calc(100vh-150px)] flex flex-col">
    <div class="mb-4">
      <h1 class="text-2xl font-bold">Общий чат</h1>
    </div>

    <div class="flex-1 bg-white rounded-lg border border-gray-200 flex flex-col overflow-hidden shadow-sm">
      <div ref="messagesContainer" class="flex-1 p-4 overflow-y-auto space-y-4 bg-gray-50 scroll-smooth">
        <div v-if="messages.length === 0" class="text-center text-gray-400 py-10">
          Сообщений пока нет. Напишите первым!
        </div>

        <div 
          v-for="msg in messages" 
          :key="msg.id"
          :class="[
            'flex flex-col max-w-[70%] p-3 rounded-lg text-sm',
            msg.isOwn ? 'bg-indigo-600 text-white self-end rounded-br-none' : 'bg-white border border-gray-200 self-start rounded-bl-none'
          ]"
        >
          <div class="flex justify-between items-baseline mb-1 opacity-80 text-[10px] font-bold">
            <span>{{ msg.username }}</span>
            <span>{{ msg.time }}</span>
          </div>
          <div class="break-words leading-relaxed">{{ msg.text }}</div>
        </div>
      </div>

      <div class="p-3 border-t border-gray-200 bg-white flex items-center gap-2">
        <Input 
          v-model="inputMessage" 
          @keyup.enter="sendMessage"
          placeholder="Напишите сообщение..." 
          :disabled="!isConnected"
          class="flex-1"
        />
        <Button @click="sendMessage" :disabled="!inputMessage.trim() || !isConnected">
          <Send class="h-4 w-4" />
        </Button>
      </div>
      
      <div class="px-3 py-1 bg-gray-100 text-[10px] flex justify-between items-center text-gray-500">
        <span>{{ isConnected ? 'Подключено' : 'Отключено' }}</span>
        <span v-if="!isConnected" class="text-orange-600 font-bold animate-pulse">Пытаемся переподключиться...</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue';
import { useAuthStore } from '~/stores/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Send } from 'lucide-vue-next';

interface Message {
  id: string;
  username: string;
  text: string;
  time: string;
  isOwn: boolean;
}

const auth = useAuthStore();
const messages = ref<Message[]>([]);
const inputMessage = ref('');
const isConnected = ref(false);
const messagesContainer = ref<HTMLElement | null>(null);

let socket: WebSocket | null = null;
let reconnectTimer: NodeJS.Timeout | null = null;

const mockMessages: Message[] = [
  { id: '1', username: 'Андрей', text: 'Привет! Каков прогресс в целом?', time: '10:00', isOwn: false },
  { id: '2', username: 'Александр', text: 'Исправляю недочёты, пока всё окей', time: '10:05', isOwn: true },
];

onMounted(() => {
  messages.value = mockMessages;
  connectWebSocket();
});

onUnmounted(() => {
  if (socket) socket.close();
  if (reconnectTimer) clearTimeout(reconnectTimer);
});

const connectWebSocket = () => {
  const wsUrl = `ws://localhost:8080/ws/chat`; 
  
  try {
    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
      console.log('WS Connected');
      isConnected.value = true;
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        addIncomingMessage(data);
      } catch (e) {
        console.error('WS Parse error', e);
      }
    };

    socket.onclose = () => {
      isConnected.value = false;
      attemptReconnect();
    };

    socket.onerror = (err) => {
      console.error('WS Error:', err);
      isConnected.value = false;
    };
  } catch (e) {
    console.error('WS Connection failed:', e);
    isConnected.value = false; 
  }
};

const attemptReconnect = () => {
  if (reconnectTimer) return;
  reconnectTimer = setTimeout(() => {
    reconnectTimer = null;
    connectWebSocket();
  }, 3000);
};

const addIncomingMessage = (data: any) => {
  messages.value.push({
    id: Date.now().toString(),
    username: data.username || 'Guest',
    text: data.text,
    time: new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}),
    isOwn: data.username === auth.user?.username
  });
  scrollToBottom();
};

const sendMessage = () => {
  if (!inputMessage.value.trim() || !isConnected.value || !socket) return;

  const messagePayload = {
    type: 'message',
    text: inputMessage.value,
    username: auth.user?.username || 'Anon',
    timestamp: Date.now()
  };

  messages.value.push({
    id: Date.now().toString(),
    username: auth.user?.username || 'Me',
    text: inputMessage.value,
    time: new Date().toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'}),
    isOwn: true
  });

  try {
    socket.send(JSON.stringify(messagePayload));
  } catch (e) {
    console.error('Send failed', e);
  }

  inputMessage.value = '';
  scrollToBottom();
};

const scrollToBottom = async () => {
  await nextTick();
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
  }
};
</script>