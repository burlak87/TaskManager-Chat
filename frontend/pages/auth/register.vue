<template>
  <div class="flex min-h-[80vh] items-center justify-center">
    <Card class="w-full max-w-md shadow-lg">
      <CardHeader>
        <CardTitle class="text-2xl text-center">Регистрация</CardTitle>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleRegister" class="space-y-4">
          <div class="space-y-2">
            <Label for="username">Имя пользователя</Label>
            <Input id="username" v-model="form.username" placeholder="user123_" required />
          </div>
          <div class="space-y-2">
            <Label for="firstname">Имя</Label>
            <Input id="firstname" v-model="form.firstname" placeholder="Иван" required />
          </div>
          <div class="space-y-2">
            <Label for="lastname">Фамилия</Label>
            <Input id="lastname" v-model="form.lastname" placeholder="Иванов" required />
          </div>
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
            {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
          </Button>
        </form>
      </CardContent>
      <CardFooter class="flex justify-center">
        <p class="text-sm text-gray-500">
          Уже есть аккаунт? <NuxtLink to="/auth/login" class="text-indigo-600 hover:underline">Войти</NuxtLink>
        </p>
      </CardFooter>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '~/stores/auth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/ui/card';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';
import { Loader2 } from 'lucide-vue-next';

const authStore = useAuthStore();
const form = reactive({
  username: '',
  firstname: '',
  lastname: '',
  email: '',
  password: '',
});
const loading = ref(false);
const errorMessage = ref('');

const router = useRouter();

const validateForm = (): string | null => {
  if (!form.username.trim()) {
    return 'Имя пользователя обязательно';
  }
  if (form.username.length < 3 || form.username.length > 50) {
    return 'Имя пользователя должно быть от 3 до 50 символов';
  }
  if (!form.firstname.trim()) {
    return 'Имя обязательно';
  }
  if (!form.lastname.trim()) {
    return 'Фамилия обязательна';
  }
  if (!form.email.trim()) {
    return 'Email обязателен';
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    return 'Некорректный формат email';
  }
  if (!form.password) {
    return 'Пароль обязателен';
  }
  if (form.password.length < 8) {
    return 'Пароль должен содержать не менее 8 символов';
  }
  if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/.test(form.password)) {
    return 'Пароль должен содержать буквы, цифры и специальные символы';
  }
  return null;
};

const handleRegister = async () => {
  const validationError = validateForm();
  if (validationError) {
    errorMessage.value = validationError;
    return;
  }

  loading.value = true;
  errorMessage.value = '';
  try {
    await authStore.register(form);
    await router.push('/dashboard');
  } catch (error: any) {
    errorMessage.value = error.message || 'Ошибка регистрации';
  } finally {
    loading.value = false;
  }
};
</script>