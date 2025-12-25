

<?php $__env->startSection('content'); ?>
<div class="px-4 py-6">
    <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold text-gray-800">Мои доски</h1>
        <a href="<?php echo e(route('boards.create')); ?>" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">
            <i class="fas fa-plus mr-2"></i>Новая доска
        </a>
    </div>

    <?php if($boards->isEmpty()): ?>
        <div class="text-center py-12">
            <i class="fas fa-columns text-6xl text-gray-300 mb-4"></i>
            <p class="text-gray-500 text-lg">Пока нет ни одной доски</p>
            <a href="<?php echo e(route('boards.create')); ?>" class="text-blue-500 hover:text-blue-600 mt-2 inline-block">
                Создайте свою первую доску →
            </a>
        </div>
    <?php else: ?>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <?php $__currentLoopData = $boards; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $board): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?>
                <div class="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
                    <div class="flex justify-between items-start mb-4">
                        <h3 class="text-xl font-semibold text-gray-800"><?php echo e($board->name); ?></h3>
                        <span class="bg-gray-100 text-gray-600 text-sm px-2 py-1 rounded">
                            <?php echo e($board->tasks->count()); ?> задач
                        </span>
                    </div>
                    
                    <div class="space-y-2 mb-4">
                        <?php
                            $todoCount = $board->tasks->where('status', 'todo')->count();
                            $inProgressCount = $board->tasks->where('status', 'in_progress')->count();
                            $doneCount = $board->tasks->where('status', 'done')->count();
                        ?>
                        
                        <div class="flex items-center">
                            <div class="w-3 h-3 bg-red-500 rounded-full mr-2"></div>
                            <span class="text-sm text-gray-600">Сделать: <?php echo e($todoCount); ?></span>
                        </div>
                        <div class="flex items-center">
                            <div class="w-3 h-3 bg-blue-500 rounded-full mr-2"></div>
                            <span class="text-sm text-gray-600">В работе: <?php echo e($inProgressCount); ?></span>
                        </div>
                        <div class="flex items-center">
                            <div class="w-3 h-3 bg-green-500 rounded-full mr-2"></div>
                            <span class="text-sm text-gray-600">Завершена: <?php echo e($doneCount); ?></span>
                        </div>
                    </div>
                    
                    <div class="flex space-x-2">
                        <a href="<?php echo e(route('boards.show', $board)); ?>" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white text-center py-2 rounded">
                            <i class="fas fa-eye mr-1"></i> Открыть
                        </a>
                        <a href="<?php echo e(route('boards.edit', $board)); ?>" class="px-3 py-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded">
                            <i class="fas fa-edit"></i>
                        </a>
                        <form action="<?php echo e(route('boards.destroy', $board)); ?>" method="POST" class="inline">
                            <?php echo csrf_field(); ?>
                            <?php echo method_field('DELETE'); ?>
                            <button type="submit" onclick="return confirm('Удалить доску?')" 
                                    class="px-3 py-2 bg-red-500 hover:bg-red-600 text-white rounded">
                                <i class="fas fa-trash"></i>
                            </button>
                        </form>
                    </div>
                </div>
            <?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?>
        </div>
    <?php endif; ?>
</div>
<?php $__env->stopSection(); ?>
<?php echo $__env->make('layouts.app', array_diff_key(get_defined_vars(), ['__data' => 1, '__path' => 1]))->render(); ?><?php /**PATH C:\Users\danka\Desktop\code\25.12\example\resources\views/boards/index.blade.php ENDPATH**/ ?>