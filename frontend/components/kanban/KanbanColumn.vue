<script setup lang="ts">
import { useKanbanStore } from '../../stores/kanban'
import KanbanCard from './KanbanCard.vue'
import TaskModal from './TaskModal.vue'
import { useKanban } from '../../composables/useKanban'

const { updateTaskStatus } = useKanban()

const props = defineProps<{ column: any }>()

const onDragEnd = (event: any) => {
    const task = props.column.tasks[event.newIndex]
    updateTaskStatus(task.id, props.column.id)
}

const store = useKanbanStore()

const createTask = (task: any) => {
    task.status = props.column.id
    task.id = Date.now()
    props.column.tasks.push(task)
}
</script>

<template>
    <article class="w-72 bg-gray-100 rounded p-3 flex flex-col">
        <h3 class="font-semibold mb-2">{{ column.title }}</h3>

        <Draggable :list="column.tasks" group="tasks" item-key="id" class="flex-1" @end="onDragEnd">
            <template #item="{ element }">
                <KanbanCard :task="element" />
            </template>
        </Draggable>

        <TaskModal @save="createTask">
            <button class="text-sm text-blue-600 mt-2">
                + Добавить задачу
            </button>
        </TaskModal>
    </article>
</template>
