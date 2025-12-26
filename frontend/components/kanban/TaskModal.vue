<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import {
    Dialog,
    DialogTrigger,
    DialogContent,
    DialogTitle
} from 'reka-ui'

const props = defineProps<{ task?: any }>()
const emit = defineEmits(['save'])

const open = ref(false)

const form = reactive({
    title: '',
    description: '',
    status: 'todo',
    assignee: '',
})

const tags = ref('')

const isEdit = computed(() => !!props.task)

watch(
    () => props.task,
    (task) => {
        if (task) {
            form.title = task.title
            form.description = task.description
            form.status = task.status
            form.assignee = task.assignee
            tags.value = task.tags?.join(', ') ?? ''
        }
    },
    { immediate: true }
)

const submit = () => {
    emit('save', {
        ...form,
        tags: tags.value.split(',').map(t => t.trim()),
    })
    open.value = false
}
</script>

<template>
    <Dialog v-model:open="open">
        <DialogTrigger as-child>
            <slot />
        </DialogTrigger>

        <DialogContent class="w-[400px]">
            <DialogTitle>
                {{ isEdit ? 'Редактировать задачу' : 'Создать задачу' }}
            </DialogTitle>

            <form class="space-y-3" @submit.prevent="submit">
                <input v-model="form.title" placeholder="Название задачи" class="input" required />

                <textarea v-model="form.description" placeholder="Описание" class="input" />

                <select v-model="form.status" class="input">
                    <option value="todo">To Do</option>
                    <option value="in_progress">In Progress</option>
                    <option value="done">Done</option>
                </select>

                <input v-model="form.assignee" placeholder="Исполнитель" class="input" />

                <input v-model="tags" placeholder="Теги через запятую" class="input" />

                <button class="btn-primary w-full">
                    Сохранить
                </button>
            </form>
        </DialogContent>
    </Dialog>
</template>