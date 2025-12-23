import { $fetch } from 'ofetch';
import { useAuthStore } from '~/stores/auth';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8888';

const api = $fetch.create({
  baseURL: `${API_BASE_URL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
  async onRequest({ options }) {
    const authStore = useAuthStore();
    if (authStore.accessToken) {
      (options.headers as any) = { ...options.headers, Authorization: `Bearer ${authStore.accessToken}` };
    }
  },
  async onResponseError({ response }) {
    if (response.status === 401) {
      const authStore = useAuthStore();
      authStore.logout();
    }
    const msg = (response._data as any)?.error || response.statusText || 'Ошибка сервера';
    throw new Error(msg);
  },
});

export const authApi = {
  login: (credentials: { email: string; password: string }) =>
    api('/auth/login', { method: 'POST', body: credentials }),
  register: (data: { username: string; firstname: string; lastname: string; email: string; password: string }) =>
    api('/auth/register', { method: 'POST', body: data }),
};

export default api;