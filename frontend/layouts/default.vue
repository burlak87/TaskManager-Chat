<template>
  <div class="min-h-screen flex flex-col bg-gray-50 text-slate-900">
    <header class="bg-white border-b border-gray-200 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <NuxtLink to="/" class="flex items-center gap-2 font-bold text-xl text-indigo-600 hover:text-indigo-700">
            <LayoutGrid class="h-5 w-5" />
            <span>TaskManager+Chat</span>
          </NuxtLink>
        </div>

        <nav v-if="auth.isAuthenticated" class="hidden md:flex items-center gap-6 text-sm font-medium text-gray-600">
          <NuxtLink to="/dashboard" class="hover:text-indigo-600">Доски</NuxtLink>
          <NuxtLink to="/chat" class="hover:text-indigo-600">Чат</NuxtLink>
        </nav>

        <div class="flex items-center gap-3">
          <div v-if="auth.isAuthenticated">
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button variant="ghost" size="icon" class="relative">
                  <Bell class="h-4 w-4" />
                  <span 
                    v-if="notification.unreadCount > 0"
                    class="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-[10px] text-white"
                  >
                    {{ notification.unreadCount }}
                  </span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" class="w-64">
                <DropdownMenuLabel>Уведомления</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <div v-if="notification.notifications.length === 0" class="p-2 text-center text-sm text-gray-500">
                  Нет новых
                </div>
                <DropdownMenuItem 
                  v-for="notif in notification.notifications" 
                  :key="notif.id"
                  @click="notification.markAsRead(notif.id)"
                  class="flex flex-col items-start gap-1 cursor-pointer"
                >
                  <span class="font-medium text-xs">{{ notif.message }}</span>
                  <span class="text-[10px] text-gray-400">{{ formatDate(notif.createdAt) }}</span>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          <div v-if="auth.isAuthenticated" class="flex items-center gap-2">
            <div class="text-right hidden sm:block">
              <div class="text-sm font-medium">{{ auth.user?.username || 'User' }}</div>
              <div class="text-xs text-gray-500">{{ auth.user?.email }}</div>
            </div>
            <Button variant="ghost" size="sm" @click="handleLogout" class="text-red-600 hover:bg-red-50 cursor-pointer">
              Выйти
            </Button>
          </div>
          <div v-else>
            <NuxtLink to="/auth/login">
              <Button class="cursor-pointer" size="sm">Войти</Button>
            </NuxtLink>
          </div>
        </div>
      </div>
    </header>

    <main class="flex-grow py-6">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <slot />
      </div>
    </main>

    <footer class="bg-white border-t border-gray-200 py-6 mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center text-sm text-gray-500">
        <p>© 2025 TaskManager+Chat. Учебный проект.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/stores/auth';
import { useNotificationStore } from '~/stores/notification';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Bell, LayoutGrid } from 'lucide-vue-next';

const auth = useAuthStore();
const notification = useNotificationStore();
const router = useRouter();

const handleLogout = () => {
  auth.logout();
  router.push('/auth/login');
};

const formatDate = (date: Date) => {
  return new Intl.DateTimeFormat('ru-RU', { 
    hour: '2-digit', 
    minute: '2-digit' 
  }).format(date);
};
</script>