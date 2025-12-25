

<?php $__env->startSection('content'); ?>
<div class="px-4 py-6">
    <div class="flex justify-between items-center mb-6">
        <div>
            <h1 class="text-3xl font-bold text-gray-800"><?php echo e($board->name); ?></h1>
            <p class="text-gray-600">Создана: <?php echo e($board->created_at->format('d.m.Y')); ?></p>
        </div>
        <div class="flex space-x-3">
            <button onclick="showCreateTaskModal()" class="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-lg">
                <i class="fas fa-plus mr-2"></i>Добавить задачу
            </button>
            <a href="<?php echo e(route('boards.edit', $board)); ?>" class="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded-lg">
                <i class="fas fa-edit mr-2"></i>Редактировать
            </a>
            <a href="<?php echo e(route('boards.index')); ?>" class="bg-gray-500 hover:bg-gray-600 text-white px-4 py-2 rounded-lg">
                <i class="fas fa-arrow-left mr-2"></i>Назад
            </a>
        </div>
    </div>

    <!-- Канбан доска -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- Сделать -->
        <div class="kanban-column">
            <div class="flex items-center mb-4">
                <div class="w-3 h-3 bg-red-500 rounded-full mr-2"></div>
                <h2 class="text-xl font-semibold text-gray-800">Сделать</h2>
                <span class="ml-2 bg-gray-200 text-gray-600 text-sm px-2 py-1 rounded">
                    <?php echo e($tasksByStatus['todo']->count()); ?>

                </span>
            </div>
            
            <?php $__currentLoopData = $tasksByStatus['todo']; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $task): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?>
                <?php echo $__env->make('partials.task-card', ['task' => $task], array_diff_key(get_defined_vars(), ['__data' => 1, '__path' => 1]))->render(); ?>
            <?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?>
            
            <?php if($tasksByStatus['todo']->isEmpty()): ?>
                <div class="text-center py-8 text-gray-400">
                    <i class="fas fa-tasks text-4xl mb-2"></i>
                    <p>Задачи отсутствуют</p>
                </div>
            <?php endif; ?>
        </div>

        <!-- В работе -->
        <div class="kanban-column">
            <div class="flex items-center mb-4">
                <div class="w-3 h-3 bg-blue-500 rounded-full mr-2"></div>
                <h2 class="text-xl font-semibold text-gray-800">В работе</h2>
                <span class="ml-2 bg-gray-200 text-gray-600 text-sm px-2 py-1 rounded">
                    <?php echo e($tasksByStatus['in_progress']->count()); ?>

                </span>
            </div>
            
            <?php $__currentLoopData = $tasksByStatus['in_progress']; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $task): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?>
                <?php echo $__env->make('partials.task-card', ['task' => $task], array_diff_key(get_defined_vars(), ['__data' => 1, '__path' => 1]))->render(); ?>
            <?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?>
            
            <?php if($tasksByStatus['in_progress']->isEmpty()): ?>
                <div class="text-center py-8 text-gray-400">
                    <i class="fas fa-cogs text-4xl mb-2"></i>
                    <p>Нет задач в работе</p>
                </div>
            <?php endif; ?>
        </div>

        <!-- Завершена -->
        <div class="kanban-column">
            <div class="flex items-center mb-4">
                <div class="w-3 h-3 bg-green-500 rounded-full mr-2"></div>
                <h2 class="text-xl font-semibold text-gray-800">Завершена</h2>
                <span class="ml-2 bg-gray-200 text-gray-600 text-sm px-2 py-1 rounded">
                    <?php echo e($tasksByStatus['done']->count()); ?>

                </span>
            </div>
            
            <?php $__currentLoopData = $tasksByStatus['done']; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $task): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?>
                <?php echo $__env->make('partials.task-card', ['task' => $task], array_diff_key(get_defined_vars(), ['__data' => 1, '__path' => 1]))->render(); ?>
            <?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?>
            
            <?php if($tasksByStatus['done']->isEmpty()): ?>
                <div class="text-center py-8 text-gray-400">
                    <i class="fas fa-check-circle text-4xl mb-2"></i>
                    <p>Нет завершенных задач</p>
                </div>
            <?php endif; ?>
        </div>
    </div>
</div>

<!-- Модальное окно создания задачи -->
<div id="createTaskModal" class="fixed inset-0 bg-black bg-opacity-50 hidden flex items-center justify-center p-4">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div class="p-6">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">Создать задачу</h2>
            
            <form id="createTaskForm" action="<?php echo e(route('tasks.store')); ?>" method="POST">
                <?php echo csrf_field(); ?>
                <input type="hidden" name="board_id" value="<?php echo e($board->id); ?>">
                
                <div class="mb-4">
                    <label for="title" class="block text-gray-700 text-sm font-medium mb-2">
                        Название задачи
                    </label>
                    <input type="text" name="title" id="title" required
                           class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                           placeholder="Что нужно сделать?">
                </div>
                
                <div class="mb-4">
                    <label for="description" class="block text-gray-700 text-sm font-medium mb-2">
                        Описание
                    </label>
                    <textarea name="description" id="description" rows="3" required
                              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                              placeholder="Подробное описание задачи"></textarea>
                </div>
                
                <div class="mb-6">
                    <label for="status" class="block text-gray-700 text-sm font-medium mb-2">
                        Статус
                    </label>
                    <select name="status" id="status" required
                            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
                        <option value="todo">Сделать</option>
                        <option value="in_progress">В работе</option>
                        <option value="done">Завершена</option>
                    </select>
                </div>
                
                <div class="flex space-x-3">
                    <button type="submit" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white py-2 rounded-lg">
                        <i class="fas fa-save mr-2"></i>Создать
                    </button>
                    <button type="button" onclick="hideCreateTaskModal()" 
                            class="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 rounded-lg">
                        <i class="fas fa-times mr-2"></i>Отмена
                    </button>
                </div>
            </form>
        </div>
    </div>
</div>

<!-- Модальное окно редактирования задачи -->
<div id="editTaskModal" class="fixed inset-0 bg-black bg-opacity-50 hidden flex items-center justify-center p-4">
    <!-- Содержимое будет заполняться JavaScript -->
</div>

<script>
function showCreateTaskModal() {
    document.getElementById('createTaskModal').classList.remove('hidden');
}

function hideCreateTaskModal() {
    document.getElementById('createTaskModal').classList.add('hidden');
}

function showEditTaskModal(taskId) {
    // Здесь можно реализовать AJAX загрузку формы редактирования
    alert('Редактирование задачи ' + taskId);
}

function showMoveTaskModal(taskId) {
    // Запрашиваем у пользователя выбор доски
    const newBoardId = prompt('Введите ID новой доски:');
    if (newBoardId) {
        const newStatus = prompt('Выберите новый статус (todo/in_progress/done):');
        if (newStatus && ['todo', 'in_progress', 'done'].includes(newStatus)) {
            if (confirm('Перенести задачу?')) {
                // Здесь можно отправить AJAX запрос
                window.location.href = `/tasks/${taskId}/move?board_id=${newBoardId}&status=${newStatus}`;
            }
        }
    }
}
</script>
<?php $__env->stopSection(); ?>
<?php echo $__env->make('layouts.app', array_diff_key(get_defined_vars(), ['__data' => 1, '__path' => 1]))->render(); ?><?php /**PATH C:\Users\danka\Desktop\code\25.12\example\resources\views/boards/show.blade.php ENDPATH**/ ?>