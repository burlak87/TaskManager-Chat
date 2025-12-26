import { defineStore } from 'pinia'

export interface Task {
  id: number
  title: string
  description?: string
  status: string
  assignee?: string
  tags?: string[]
}

export interface Column {
  id: string
  title: string
  tasks: Task[]
}

export let useKanbanStore = defineStore('kanban', {
  state: () => ({
    columns: [] as Column[],
    loading: false,
  }),

  actions: {
    setColumns(columns: Column[]) {
      this.columns = columns
    },

    moveTask(taskId: number, from: string, to: string) {
      let fromCol = this.columns.find(c => c.id === from)
      let toCol = this.columns.find(c => c.id === to)
      if (!fromCol || !toCol) return

      let idx = fromCol.tasks.findIndex(t => t.id === taskId)
      let task = fromCol.tasks.splice(idx, 1)[0]
      task.status = to
      toCol.tasks.push(task)
    },
  },
})
