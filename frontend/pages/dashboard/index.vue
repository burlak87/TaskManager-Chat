<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold">Ваши проекты</h1>
      <Dialog v-model:open="showCreateDialog">
        <DialogTrigger as-child>
          <Button class="cursor-pointer">
            <Plus class="mr-2 h-4 w-4" /> Новый проект
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Создать новый проект</DialogTitle>
            <DialogDescription>
              Назовите вашу доску задач
            </DialogDescription>
          </DialogHeader>
          <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="title" class="text-right">Название</Label>
              <Input id="title" v-model="newBoardTitle" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="description" class="text-right">Описание</Label>
              <Textarea id="description" v-model="newBoardDescription" class="col-span-3" />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" @click="showCreateDialog = false">Отмена</Button>
            <Button @click="createBoard">Создать</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>

    <div v-if="boardStore.loading" class="flex justify-center py-10">
      <Loader2 class="h-8 w-8 animate-spin text-indigo-600" />
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <NuxtLink 
        v-for="board in boardStore.boards" 
        :key="board.id"
        :to="`/dashboard/${board.id}`"
        class="group block p-6 bg-white rounded-lg border border-gray-200 shadow-sm hover:shadow-md hover:border-indigo-500 transition-all"
      >
        <div class="flex justify-between items-start">
          <div>
            <h2 class="text-xl font-semibold text-gray-900 group-hover:text-indigo-600">{{ board.title }}</h2>
            <p class="mt-2 text-sm text-gray-500">{{ board.description || 'Без описания' }}</p>
          </div>
          <ArrowRight class="h-5 w-5 text-gray-400 group-hover:text-indigo-600" />
        </div>
      </NuxtLink>

      <div v-if="boardStore.boards.length === 0" class="col-span-full text-center py-12 text-gray-500">
        <p>У вас еще нет проектов. Создайте первый!</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useBoardStore } from '~/stores/board';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter
} from '@/components/ui/dialog';
import { Plus, ArrowRight, Loader2 } from 'lucide-vue-next';

const boardStore = useBoardStore();

const showCreateDialog = ref(false);
const newBoardTitle = ref('');
const newBoardDescription = ref('');

onMounted(() => {
  if (import.meta.client) {
    boardStore.fetchBoards();
  }
});

const createBoard = async () => {
  if (!newBoardTitle.value.trim()) return;

  try {
    await boardStore.createBoard(newBoardTitle.value, newBoardDescription.value);
    newBoardTitle.value = '';
    newBoardDescription.value = '';
    showCreateDialog.value = false;
  } catch (e) {
    console.error('Error creating board:', e);
  }
};
</script>