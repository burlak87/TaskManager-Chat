<template>
  <div class="flex min-h-[80vh] items-center justify-center">
    <Card class="w-full max-w-md shadow-lg">
      <CardHeader>
        <CardTitle class="text-2xl text-center">Вход</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div class="space-y-2">
            <Label for="email">Email</Label>
            <Input id="email" type="email" v-model="form.email" placeholder="user@example.com" required />
          </div>
          <div class="space-y-2">
            <Label for="password">Пароль</Label>
            <Input id="password" type="password" v-model="form.password" required />
          </div>
          <Alert v-if="errorMessage" variant="destructive">
            <AlertTitle>Ошибка</AlertTitle>
            <AlertDescription>{{ errorMessage }}</AlertDescription>
          </Alert>
          <Button type="submit" class="w-full cursor-pointer" :disabled="loading">
            <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
            {{ loading ? 'Вход...' : 'Войти' }}
          </Button>
        </form>
        <div class="mt-6 text-center text-sm">
          Нет аккаунта?
          <NuxtLink to="/auth/register" class="text-indigo-600 hover:underline font-medium">
            Зарегистрироваться
          </NuxtLink>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '~/stores/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';
import { Loader2 } from 'lucide-vue-next';

const authStore = useAuthStore();
const form = reactive({ email: '', password: '' });
const loading = ref(false);
const errorMessage = ref('');

const router = useRouter();

const handleLogin = async () => {
  loading.value = true;
  errorMessage.value = '';
  try {
    await authStore.login(form);
    await router.push('/dashboard');
  } catch (error: any) {
    errorMessage.value = error.message || 'Не удалось войти';
  } finally {
    loading.value = false;
  }
};
</script>