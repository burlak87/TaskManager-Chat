import { defineStore } from 'pinia';
import { boardsApi } from '~/services/api';
import { useNotificationStore } from '~/stores/notification';

export interface Column {
  id: string;
  title: string;
  order: number;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  status: string;
  assignee?: string;
  tags?: string[];
  createdAt?: string;
}

export const useBoardStore = defineStore('board', {
  state: () => ({
    boards: [] as any[],
    currentBoard: null as any,
    columns: [
      { id: 'todo', title: 'To Do', order: 1 },
      { id: 'in-progress', title: 'In Progress', order: 2 },
      { id: 'done', title: 'Done', order: 3 },
    ] as Column[],
    tasks: [] as Task[],
    loading: false,
  }),

  actions: {
    async fetchBoards() {
      this.loading = true;
      try {
        const response = await boardsApi.getBoards();
        this.boards = response;
      } catch (e) {
        console.error('Error fetching boards:', e);
      } finally {
        this.loading = false;
      }
    },

    async createBoard(title: string, description?: string) {
      try {
        const response = await boardsApi.createBoard({ title, description });
        this.boards.push(response);
        return response;
      } catch (e) {
        console.error('Error creating board:', e);
        throw e;
      }
    },

    async fetchTasks(boardId: string) {
      this.loading = true;
      try {
        const response = await boardsApi.getTasks(boardId);
        this.tasks = response;
      } catch (e) {
        console.error('Error fetching tasks:', e);
      } finally {
        this.loading = false;
      }
    },

    async addTask(task: Omit<Task, 'id'> & { boardId: string }) {
      try {
        const { boardId, ...taskData } = task;
        const response = await boardsApi.createTask(boardId, taskData);
        this.tasks.push(response);

        const notifStore = useNotificationStore();
        notifStore.addNotification({
          type: 'info',
          message: `Задача "${task.title}" создана`
        });

        return response;
      } catch (e) {
        console.error('Error adding task:', e);
        throw e;
      }
    },

    async updateTaskStatus(taskId: string, newStatus: string, boardId: string) {
      try {
        const task = this.tasks.find(t => t.id === taskId);
        if (!task) return;

        const response = await boardsApi.updateTask(boardId, taskId, {
          status: newStatus
        });

        if (response) {
          Object.assign(task, response);
        }

        const notifStore = useNotificationStore();
        notifStore.addNotification({
          type: 'success',
          message: `Задача перемещена в ${newStatus}`
        });
      } catch (e) {
        console.error('Error updating task status:', e);
        throw e;
      }
    }
  }
});