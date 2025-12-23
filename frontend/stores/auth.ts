import { defineStore } from 'pinia';
import { authApi } from '~/services/api';

interface User {
  id?: number;
  username: string;
  email: string;
  firstname?: string;
  lastname?: string;
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    accessToken: typeof window !== 'undefined' ? localStorage.getItem('access_token') : null,
    loading: false,
  }),
  getters: {
    isAuthenticated: (state) => !!state.accessToken,
  },
  actions: {
    setToken(token: string) {
      this.accessToken = token;
      localStorage.setItem('access_token', token);
    },
    setUser(user: User) {
      this.user = user;
    },
    async login(credentials: { email: string; password: string }) {
      this.loading = true;
      try {
        const data = await authApi.login(credentials);
        this.setToken(data.access_token);
        this.setUser({
          email: credentials.email,
          username: credentials.email.split('@')[0],
        });
        return true;
      } catch (error: any) {
        throw new Error(error.message);
      } finally {
        this.loading = false;
      }
    },
    async register(form: { username: string; firstname: string; lastname: string; email: string; password: string }) {
      this.loading = true;
      try {
        await authApi.register(form);
        await this.login({ email: form.email, password: form.password });
      } catch (error: any) {
        throw new Error(error.message);
      } finally {
        this.loading = false;
      }
    },
    logout() {
      this.user = null;
      this.accessToken = null;
      localStorage.removeItem('access_token');
      navigateTo('/auth/login');
    },
  },
});