<script setup lang="ts">
import { useDebounceFn } from '@vueuse/core'
import { ref, watch } from 'vue'

const emit = defineEmits(['change'])

let search = ref('')
let status = ref('')

const emitChange = useDebounceFn(() => {
  emit('change', { search: search.value, status: status.value })
}, 300)

watch([search, status], emitChange)
</script>

<template>
  <article class="flex gap-3 mb-4">
    <input v-model="search" placeholder="Поиск..." class="input w-64" />
    <select v-model="status" class="input">
      <option value="">Все статусы</option>
      <option value="todo">To Do</option>
      <option value="in_progress">In Progress</option>
      <option value="done">Done</option>
    </select>
  </article>
</template>
