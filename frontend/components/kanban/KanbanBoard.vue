<script setup lang="ts">
import { ref, computed } from 'vue'
import { useKanban } from '../../composables/useKanban'
import KanbanColumn from './KanbanColumn.vue'
import TaskFilters from './TaskFilters.vue'

const { store } = useKanban()

const filters = ref({ search: '', status: '' })

const filteredColumns = computed(() => {
    return store.columns.map(col => ({
        ...col,
        tasks: col.tasks.filter(t => {
            const matchTitle = t.title
                .toLowerCase()
                .includes(filters.value.search.toLowerCase())

            const matchStatus = filters.value.status
                ? t.status === filters.value.status
                : true

            return matchTitle && matchStatus
        }),
    }))
})
</script>

<template>
    <section>
        <TaskFilters @change="filters = $event" />

        <article class="flex gap-4 overflow-x-auto py-4 px-2 min-h-[70vh] items-start">
            <KanbanColumn v-for="column in filteredColumns" :key="column.id" :column="column" />

            <article v-if="!filteredColumns.length" class="text-gray-400 text-sm italic">
                Нет задач по выбранным фильтрам
            </article>
        </article>

    </section>
</template>
