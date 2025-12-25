<script setup lang="ts">
import KanbanBoard from '~/components/kanban/KanbanBoard.vue'
import { useKanban } from '~/composables/useKanban'
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'


const route = useRoute()
const { fetchBoard } = useKanban()

onMounted(() => {
    fetchBoard(route.params.id as string)
})
</script>

<template>
    <KanbanBoard />
</template>