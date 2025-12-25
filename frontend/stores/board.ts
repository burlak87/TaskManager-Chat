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
  tags?: string[];
  createdAt?: string;
}

export const useBoardStore = defineStore('board', {
  state: () => ({
    boards: [] as any[],
    currentBoard: null as any,
    columns: [] as Column[],
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

        const defaultColumns = [
          { title: 'Сделать' },
          { title: 'В работе' },
          { title: 'Готово' }
        ];

        for (const col of defaultColumns) {
          await boardsApi.createColumn(response.id, { title: col.title });
        }

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
        this.tasks = response.map((task: any) => ({
          ...task,
          status: task.column_id?.toString()
        }));
      } catch (e) {
        console.error('Error fetching tasks:', e);
      } finally {
        this.loading = false;
      }
    },

    async addTask(task: { title: string; description?: string; column_id: number | string; boardId: string }) {
      try {
        const columnId = typeof task.column_id === 'string' ? parseInt(task.column_id) : task.column_id;

        const response = await boardsApi.createTask(task.boardId, {
          title: task.title,
          description: task.description,
          column_id: columnId
        });
        const taskWithStatus = {
          ...response,
          status: response.column_id?.toString()
        };
        this.tasks.push(taskWithStatus);

        const notifStore = useNotificationStore();
        notifStore.addNotification({
          type: 'info',
          message: `Задача "${task.title}" создана`
        });

        return taskWithStatus;
      } catch (e) {
        console.error('Error adding task:', e);
        throw e;
      }
    },

    async updateTask(taskId: string, boardId: string, updateData: { title?: string; description?: string; status?: string }) {
      try {
        const task = this.tasks.find(t => t.id === taskId);
        if (!task) return;

        const requestData: any = { ...updateData };
        if (updateData.status) {
          requestData.column_id = parseInt(updateData.status);
          delete requestData.status;
        }

        const response = await boardsApi.updateTask(boardId, taskId, requestData);

        if (response) {
          const taskWithStatus = {
            ...response,
            status: response.column_id?.toString() || task.status
          };
          Object.assign(task, taskWithStatus);
        }

        const notifStore = useNotificationStore();
        notifStore.addNotification({
          type: 'success',
          message: 'Задача обновлена'
        });

        return response;
      } catch (e) {
        console.error('Error updating task:', e);
        throw e;
      }
    },

    async updateTaskStatus(taskId: string, newStatus: string, boardId: string) {
      try {
        const task = this.tasks.find(t => t.id === taskId);
        if (!task) return;

        const response = await boardsApi.updateTask(boardId, taskId, {
          column_id: parseInt(newStatus)
        });

        if (response) {
          const taskWithStatus = {
            ...response,
            status: response.column_id?.toString()
          };
          Object.assign(task, taskWithStatus);
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
    },

    async fetchColumns(boardId: string) {
      try {
        const response = await boardsApi.getColumns(boardId);
        this.columns = response.map((col: any) => ({
          id: col.id?.toString() || col.id,
          title: col.title,
          order: col.position || col.order
        }));
      } catch (e) {
        console.error('Error fetching columns:', e);
        throw e;
      }
    },

    async addColumn(boardId: string, title: string, position?: number) {
      try {
        const response = await boardsApi.createColumn(boardId, { title, position });
        this.columns.push(response);
        return response;
      } catch (e) {
        console.error('Error adding column:', e);
        throw e;
      }
    },

    async updateColumn(boardId: string, columnId: string, title: string) {
      try {
        const response = await boardsApi.updateColumn(boardId, columnId, { title });
        const column = this.columns.find(c => c.id === columnId);
        if (column) {
          Object.assign(column, response);
        }
        return response;
      } catch (e) {
        console.error('Error updating column:', e);
        throw e;
      }
    },

    async deleteTask(taskId: string, boardId: string) {
      try {
        await boardsApi.deleteTask(boardId, taskId);
        this.tasks = this.tasks.filter(t => t.id !== taskId);

        const notifStore = useNotificationStore();
        notifStore.addNotification({
          type: 'info',
          message: 'Задача удалена'
        });
      } catch (e) {
        console.error('Error deleting task:', e);
        throw e;
      }
    },

    async deleteColumn(boardId: string, columnId: string) {
      try {
        await boardsApi.deleteColumn(boardId, columnId);
        this.columns = this.columns.filter(c => c.id !== columnId);
      } catch (e) {
        console.error('Error deleting column:', e);
        throw e;
      }
    }
  }
});