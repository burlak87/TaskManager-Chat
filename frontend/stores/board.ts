import { defineStore } from 'pinia';
import { boardsApi } from '~/services/api';

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
        this.boards = [
          { id: '1', title: 'Backend MVP', description: 'API для чата и задач' },
          { id: '2', title: 'Frontend Layout', description: 'Nuxt и Tailwind' },
        ];
      } catch (e) {
        console.error('Error fetching boards:', e);
      } finally {
        this.loading = false;
      }
    },

    async createBoard(title: string) {
      try {
        const newBoard = { id: Date.now().toString(), title };
        this.boards.push(newBoard);
        return newBoard;
      } catch (e) {
        console.error('Error creating board:', e);
        throw e;
      }
    },

    async fetchTasks(boardId: string) {
      this.loading = true;
      try {
        this.tasks = [
          { id: '1', title: 'Настроить Nuxt', description: 'Сделать шапку и роуты', status: 'done', assignee: 'Саня', createdAt: '2025-12-08' },
          { id: '2', title: 'Сделать API', description: 'Аксис и интерцепторы', status: 'in-progress', assignee: 'Федор', createdAt: '2025-12-09' },
          { id: '3', title: 'Драг-н-дроп', description: 'Реализовать перетаскивание', status: 'todo', assignee: 'Андрей', createdAt: '2025-12-10' },
        ];
      } catch (e) {
        console.error('Error fetching tasks:', e);
      } finally {
        this.loading = false;
      }
    },

    async addTask(task: Omit<Task, 'id'>) {
      try {
        const newTask = { ...task, id: Date.now().toString() };
        this.tasks.push(newTask);
        
        const notifStore = useNotificationStore();
        notifStore.addNotification({
          type: 'info',
          message: `Задача "${task.title}" создана`
        });
        
        return newTask;
      } catch (e) {
        console.error('Error adding task:', e);
        throw e;
      }
    },

    async updateTaskStatus(taskId: string, newStatus: string) {
      const taskIndex = this.tasks.findIndex(t => t.id === taskId);
      if (taskIndex !== -1) {
        const oldStatus = this.tasks[taskIndex].status;
        this.tasks[taskIndex].status = newStatus;

        try {
          const notifStore = useNotificationStore();
          notifStore.addNotification({
            type: 'success',
            message: `Задача перемещена в ${newStatus}`
          });
        } catch (e) {
          this.tasks[taskIndex].status = oldStatus;
          console.error('Error updating status:', e);
        }
      }
    }
  }
});