<div class="task-card bg-white rounded-lg shadow p-4 status-<?php echo e($task->status); ?>">
    <div class="flex justify-between items-start mb-2">
        <h3 class="font-semibold text-gray-800"><?php echo e($task->title); ?></h3>
        <div class="flex space-x-1">
            <button onclick="showEditTaskModal(<?php echo e($task->id); ?>)" 
                    class="text-yellow-500 hover:text-yellow-600">
                <i class="fas fa-edit"></i>
            </button>
            <button onclick="showMoveTaskModal(<?php echo e($task->id); ?>)"
                    class="text-blue-500 hover:text-blue-600">
                <i class="fas fa-exchange-alt"></i>
            </button>
            <form action="<?php echo e(route('tasks.destroy', $task)); ?>" method="POST" class="inline">
                <?php echo csrf_field(); ?>
                <?php echo method_field('DELETE'); ?>
                <button type="submit" onclick="return confirm('Удалить задачу?')"
                        class="text-red-500 hover:text-red-600">
                    <i class="fas fa-trash"></i>
                </button>
            </form>
        </div>
    </div>
    
    <p class="text-gray-600 text-sm mb-3"><?php echo e(Str::limit($task->description, 100)); ?></p>
    
    <div class="flex justify-between items-center text-xs text-gray-500">
        <span>ID: <?php echo e($task->id); ?></span>
        <span><?php echo e($task->created_at->format('d.m.Y')); ?></span>
    </div>
</div><?php /**PATH C:\Users\danka\Desktop\code\25.12\example\resources\views/partials/task-card.blade.php ENDPATH**/ ?>