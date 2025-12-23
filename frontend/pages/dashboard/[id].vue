<template>
  <div class="flex flex-col h-[calc(100vh-150px)]">
    <div class="flex justify-between items-center mb-4 shrink-0">
      <div class="flex items-center gap-2">
        <Button variant="ghost" size="sm" @click="navigateTo('/dashboard')">
          <ArrowLeft class="h-4 w-4 mr-1" /> Назад
        </Button>
        <h1 class="text-2xl font-bold">{{ currentBoardTitle }}</h1>
      </div>

      <div class="flex bg-gray-100 p-1 rounded-lg">
        <button
          @click="currentTab = 'board'"
          :class="tabClass('board')"
        >
          Доска
        </button>
        <button
          @click="currentTab = 'chat'"
          :class="tabClass('chat')"
        >
          Чат
          <span v-if="unreadChatCount > 0" class="ml-1 inline-block w-2 h-2 bg-red-500 rounded-full"></span>
        </button>
      </div>
    </div>

    <div v-if="boardStore.loading" class="flex-1 flex items-center justify-center">
      <Loader2 class="h-8 w-8 animate-spin text-indigo-600" />
    </div>

    <div v-else-if="currentTab === 'board'" class="flex-1 overflow-x-auto bg-gray-50 rounded-lg border border-gray-200 p-4 flex gap-4">
      <div
        v-for="col in boardStore.columns"
        :key="col.id"
        class="min-w-[280px] w-[280px] bg-gray-100 rounded-md flex flex-col max-h-full"
      >
        <div class="p-3 font-semibold text-gray-700 flex justify-between border-b border-gray-200">
          {{ col.title }}
          <span class="bg-gray-200 text-gray-600 text-xs px-2 py-0.5 rounded-full">
            {{ boardStore.tasks.filter(t => t.status === col.id).length }}
          </span>
        </div>

        <div class="p-2 flex flex-col gap-2 overflow-y-auto flex-1">
          <div
            v-for="task in boardStore.tasks.filter(t => t.status === col.id)"
            :key="task.id"
            class="bg-white p-3 rounded shadow-sm border border-gray-200 text-sm group"
          >
            <div class="font-medium">{{ task.title }}</div>
            <div class="text-gray-500 mt-1 truncate">{{ task.description }}</div>
            <div v-if="task.assignee" class="mt-2 flex items-center gap-1 text-xs text-gray-400">
              <div class="w-4 h-4 bg-indigo-100 rounded-full flex items-center justify-center text-indigo-700 font-bold">
                {{ task.assignee[0].toUpperCase() }}
              </div>
              {{ task.assignee }}
            </div>
          </div>

          <Button variant="outline" size="sm" class="w-full mt-1" @click="openAddTaskModal(col.id)">
            + Задача
          </Button>
        </div>
      </div>
    </div>

    <div v-else class="flex-1 flex flex-col bg-white rounded-lg border border-gray-200 overflow-hidden">
      <div ref="messagesContainer" class="flex-1 overflow-y-auto p-4 space-y-4 bg-gray-50">
        <div v-if="messages.length === 0" class="text-center text-gray-400 py-10">
          Сообщений пока нет. Напишите первым!
        </div>
        <div
          v-for="msg in messages"
          :key="msg.id"
          :class="[
            'max-w-[70%] p-3 rounded-lg text-sm',
            msg.isOwn ? 'bg-indigo-600 text-white self-end rounded-br-none' : 'bg-white border border-gray-200 self-start rounded-bl-none'
          ]"
        >
          <div class="text-xs opacity-75 mb-1">{{ msg.username }} • {{ msg.time }}</div>
          <div>{{ msg.content }}</div>
        </div>
      </div>

      <div class="p-3 border-t border-gray-200 bg-white flex gap-2">
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

      <div class="px-3 py-1 text-xs text-center text-gray-500">
        {{ isConnected ? 'Подключено к чату' : 'Подключение...' }}
      </div>
    </div>

    <Dialog v-model:open="showTaskModal">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Новая задача</DialogTitle>
        </DialogHeader>
        <div class="grid gap-4 py-4">
          <div class="space-y-2">
            <Label>Заголовок</Label>
            <Input v-model="newTask.title" placeholder="Что нужно сделать?" />
          </div>
          <div class="space-y-2">
            <Label>Описание</Label>
            <Textarea v-model="newTask.description" placeholder="Детали задачи..." />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="showTaskModal = false">Отмена</Button>
          <Button @click="addTask">Создать</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '~/stores/auth';
import { useBoardStore } from '~/stores/board';
import { useNotificationStore } from '~/stores/notification';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter
} from '@/components/ui/dialog';
import { ArrowLeft, Send, Loader2 } from 'lucide-vue-next';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const boardStore = useBoardStore();
const notifStore = useNotificationStore();

const boardId = computed(() => route.params.id as string);
const currentBoardTitle = computed(() => {
  const board = boardStore.boards.find(b => b.id === boardId.value);
  return board?.title || 'Доска';
});

const currentTab = ref<'board' | 'chat'>('board');
const unreadChatCount = ref(0);

const messages = ref<Array<{ id: string; username: string; content: string; time: string; isOwn: boolean }>>([]);
const inputMessage = ref('');
const isConnected = ref(false);
const messagesContainer = ref<HTMLElement | null>(null);
let socket: WebSocket | null = null;

const showTaskModal = ref(false);
const currentColumnId = ref('');
const newTask = reactive({
  title: '',
  description: ''
});

const tabClass = (tab: 'board' | 'chat') => [
  'px-4 py-1.5 rounded-md text-sm font-medium transition-colors',
  currentTab.value === tab ? 'bg-white shadow text-indigo-600' : 'text-gray-600 hover:text-gray-900'
];

const connectWebSocket = () => {
  if (!boardId.value || !authStore.accessToken) return;

  const wsUrl = `ws://localhost:8888/ws/chat?board_id=${boardId.value}`;
  socket = new WebSocket(wsUrl);

  socket.onopen = () => {
    isConnected.value = true;
    console.log('WebSocket подключён к доске', boardId.value);
  };

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if (data.type === 'message' && data.payload) {
        const p = data.payload;
        messages.value.push({
          id: p.id || Date.now().toString(),
          username: p.username || 'Аноним',
          content: p.content,
          time: new Date(p.created_at || Date.now()).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' }),
          isOwn: p.username === authStore.user?.username
        });

        if (currentTab.value !== 'chat') {
          unreadChatCount.value++;
        }

        scrollToBottom();
      }
    } catch (e) {
      console.error('Ошибка обработки сообщения WS:', e);
    }
  };

  socket.onclose = () => {
    isConnected.value = false;
    console.log('WebSocket закрыт, пытаемся переподключиться...');
    setTimeout(connectWebSocket, 3000);
  };

  socket.onerror = (err) => {
    console.error('Ошибка WebSocket:', err);
  };
};

const sendMessage = () => {
  if (!socket || socket.readyState !== WebSocket.OPEN || !inputMessage.value.trim()) return;

  const payload = {
    type: 'message',
    payload: {
      content: inputMessage.value.trim(),
      board_id: boardId.value
    }
  };

  socket.send(JSON.stringify(payload));

  messages.value.push({
    id: Date.now().toString(),
    username: authStore.user?.username || 'Я',
    content: inputMessage.value.trim(),
    time: new Date().toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' }),
    isOwn: true
  });

  inputMessage.value = '';
  scrollToBottom();
};

const scrollToBottom = async () => {
  await nextTick();
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
  }
};

watch(currentTab, (newTab) => {
  if (newTab === 'chat') {
    connectWebSocket();
    unreadChatCount.value = 0;
    nextTick(scrollToBottom);
  } else {
    if (socket) {
      socket.close();
      socket = null;
    }
  }
});

onMounted(async () => {
  if (boardId.value) {
    await boardStore.fetchTasks(boardId.value);
  }
});

onUnmounted(() => {
  if (socket) {
    socket.close();
  }
});

const openAddTaskModal = (colId: string) => {
  currentColumnId.value = colId;
  newTask.title = '';
  newTask.description = '';
  showTaskModal.value = true;
};

const addTask = async () => {
  if (!newTask.title.trim()) return;

  try {
    await boardStore.addTask({
      title: newTask.title,
      description: newTask.description,
      status: currentColumnId.value,
      boardId: boardId.value
    });
    showTaskModal.value = false;
  } catch (e) {
    console.error('Ошибка создания задачи:', e);
  }
};
</script>