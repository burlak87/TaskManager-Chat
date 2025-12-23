import { defineStore } from 'pinia';

export interface Notification {
  id: string;
  type: 'info' | 'success' | 'warning' | 'error';
  message: string;
  read: boolean;
  createdAt: Date;
}

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [] as Notification[],
  }),

  getters: {
    unreadCount: (state) => state.notifications.filter(n => !n.read).length,
    unreadNotifications: (state) => state.notifications.filter(n => !n.read),
  },

  actions: {
    addNotification(notification: Omit<Notification, 'id' | 'createdAt' | 'read'>) {
      const newNotification = {
        ...notification,
        id: Date.now().toString(),
        createdAt: new Date(),
        read: false,
      };
      this.notifications.unshift(newNotification);
      
      console.log('New Notification:', newNotification.message);
    },

    markAsRead(id: string) {
      const notif = this.notifications.find(n => n.id === id);
      if (notif) notif.read = true;
    },

    markAllAsRead() {
      this.notifications.forEach(n => n.read = true);
    }
  }
});