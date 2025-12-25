<template>
  <div class="flex flex-col items-center justify-center min-h-[60vh]">
    <h1 class="text-4xl font-bold mb-4">TaskManager+Chat</h1>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useAuthStore } from '~/stores/auth';

const authStore = useAuthStore();

onMounted(async () => {
  if (import.meta.client) {
    await new Promise(resolve => setTimeout(resolve, 0));
    if (authStore.isAuthenticated) {
      await navigateTo('/dashboard');
    } else {
      await navigateTo('/auth/login');
    }
  }
});
</script>