import axios from 'axios';
import { useAuthStore } from '~/stores/auth';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: `${API_BASE_URL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore();
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore();
      authStore.logout();
    }

    const msg = error.response?.data?.error || error.response?.statusText || 'Ошибка сервера';
    return Promise.reject(new Error(msg));
  }
);

export const authApi = {
  login: (credentials: { email: string; password: string }) =>
    api.post('/auth/login', credentials).then(res => res.data),
  register: (data: { username: string; firstname: string; lastname: string; email: string; password: string }) =>
    api.post('/auth/register', data).then(res => res.data),
  getProfile: () =>
    api.get('/profile').then(res => res.data),
};

export const boardsApi = {
  getBoards: () => api.get('/boards').then(res => res.data),
  createBoard: (data: { title: string; description?: string }) =>
    api.post('/boards', data).then(res => res.data),
  getBoard: (id: string) => api.get(`/boards/${id}`).then(res => res.data),
  updateBoard: (id: string, data: { title?: string; description?: string }) =>
    api.put(`/boards/${id}`, data).then(res => res.data),
  deleteBoard: (id: string) =>
    api.delete(`/boards/${id}`).then(res => res.data),

  getColumns: (boardId: string) => api.get(`/boards/${boardId}/columns`).then(res => res.data),
  createColumn: (boardId: string, data: { title: string; position?: number }) =>
    api.post(`/boards/${boardId}/columns`, data).then(res => res.data),
  updateColumn: (boardId: string, columnId: string, data: { title?: string; position?: number }) =>
    api.put(`/boards/${boardId}/columns/${columnId}`, data).then(res => res.data),
  deleteColumn: (boardId: string, columnId: string) =>
    api.delete(`/boards/${boardId}/columns/${columnId}`).then(res => res.data),

  getTasks: (boardId: string) =>
    api.get(`/boards/${boardId}/tasks`).then(res => res.data),
  createTask: (boardId: string, data: { title: string; description?: string; column_id?: number | string }) =>
    api.post(`/boards/${boardId}/tasks`, data).then(res => res.data),
  updateTask: (boardId: string, taskId: string, data: { title?: string; description?: string; status?: string; column_id?: number | string }) =>
    api.put(`/boards/${boardId}/tasks/${taskId}`, data).then(res => res.data),
  deleteTask: (boardId: string, taskId: string) =>
    api.delete(`/boards/${boardId}/tasks/${taskId}`).then(res => res.data),
};

export default api;