

<?php $__env->startSection('content'); ?>
<div class="max-w-md mx-auto bg-white rounded-lg shadow-md p-6">
    <h1 class="text-2xl font-bold text-gray-800 mb-6">Создать новую доску</h1>
    
    <form action="<?php echo e(route('boards.store')); ?>" method="POST">
        <?php echo csrf_field(); ?>
        
        <div class="mb-4">
            <label for="name" class="block text-gray-700 text-sm font-medium mb-2">
                Название доски
            </label>
            <input type="text" name="name" id="name" required
                   class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                   placeholder="Например: Проект разработки">
            <?php $__errorArgs = ['name'];
$__bag = $errors->getBag($__errorArgs[1] ?? 'default');
if ($__bag->has($__errorArgs[0])) :
if (isset($message)) { $__messageOriginal = $message; }
$message = $__bag->first($__errorArgs[0]); ?>
                <p class="text-red-500 text-sm mt-1"><?php echo e($message); ?></p>
            <?php unset($message);
if (isset($__messageOriginal)) { $message = $__messageOriginal; }
endif;
unset($__errorArgs, $__bag); ?>
        </div>
        
        <div class="flex space-x-3">
            <button type="submit" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white py-2 rounded-lg">
                <i class="fas fa-save mr-2"></i>Создать доску
            </button>
            <a href="<?php echo e(route('boards.index')); ?>" class="flex-1 bg-gray-500 hover:bg-gray-600 text-white py-2 rounded-lg text-center">
                <i class="fas fa-times mr-2"></i>Отмена
            </a>
        </div>
    </form>
</div>
<?php $__env->stopSection(); ?>
<?php echo $__env->make('layouts.app', array_diff_key(get_defined_vars(), ['__data' => 1, '__path' => 1]))->render(); ?><?php /**PATH C:\Users\danka\Desktop\code\25.12\example\resources\views/boards/create.blade.php ENDPATH**/ ?>