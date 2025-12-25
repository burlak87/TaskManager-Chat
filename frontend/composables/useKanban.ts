import { useKanbanStore } from '../stores/kanban'

export const useKanban = () => {
    const store = useKanbanStore()
    const { $fetch } = useNuxtApp()

    const fetchBoard = async (boardId: string) => {
        store.loading = true

        const data = await $fetch<{ columns: any[] }>(
            `/api/boards/${boardId}/tasks`
        )


        store.setColumns(data.columns)
        store.loading = false
    }

    const updateTaskStatus = async (taskId: number, status: string) => {
    const previous = JSON.parse(JSON.stringify(store.columns))

    try {
        await $fetch(`/api/tasks/${taskId}`, {
        method: 'PATCH',
        body: { status },
        })
    } catch (e) {
        store.setColumns(previous)
        throw e
    }
    }

    return {
        store,
        fetchBoard,
        updateTaskStatus,
    }
}
