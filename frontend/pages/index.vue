<template>
  <div class="flex flex-col items-center justify-center min-h-[60vh]">
    <h1 class="text-4xl font-bold mb-4">TaskManager+Chat</h1>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useAuthStore } from '~/stores/auth';

const authStore = useAuthStore();

onMounted(() => {
    if (authStore.isAuthenticated) {
      navigateTo('/dashboard');
    } else {
      navigateTo('/auth/login');
    }
  }
);
</script>