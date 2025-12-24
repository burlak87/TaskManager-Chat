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
        const response = await authApi.login(credentials);
        if (response.access_token) {
          this.setToken(response.access_token);
        }

        if (response.user) {
          this.setUser(response.user);
        } else {
          this.setUser({
            email: credentials.email,
            username: credentials.email.split('@')[0] || credentials.email,
          });
        }

        return true;
      } catch (error: any) {
        throw new Error(error.message || 'Ошибка входа');
      } finally {
        this.loading = false;
      }
    },
    async register(form: { username: string; firstname: string; lastname: string; email: string; password: string }) {
      this.loading = true;
      try {
        const response = await authApi.register(form);
        if (response.access_token) {
          this.setToken(response.access_token);
          if (response.user) {
            this.setUser(response.user);
          }
        } else {
          await this.login({ email: form.email, password: form.password });
        }
        return true;
      } catch (error: any) {
        throw new Error(error.message || 'Ошибка регистрации');
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