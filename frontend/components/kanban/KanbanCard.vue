<script setup lang="ts">
import TaskModal from './TaskModal.vue'

const props = defineProps<{ task: any }>()

const updateTask = (updated: any) => {
    Object.assign(props.task, updated)
}
</script>

<template>
    <TaskModal :task="task" @save="updateTask">
        <article class="bg-white rounded shadow p-2 mb-2 cursor-pointer">
            <article class="h-1 rounded mb-1" :class="{
                'bg-gray-400': task.status === 'todo',
                'bg-yellow-400': task.status === 'in_progress',
                'bg-green-500': task.status === 'done',
            }" />

            <article class="font-medium">{{ task.title }}</article>

            <article class="text-xs text-gray-500 mt-1">
                {{ task.assignee || 'Не назначен' }}
            </article>

            <article class="flex gap-1 mt-1">
                <span v-for="tag in task.tags" :key="tag" class="text-xs bg-blue-100 px-1 rounded">
                    {{ tag }}
                </span>
            </article>
        </article>
    </TaskModal>
</template>
