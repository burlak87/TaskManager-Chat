<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Канбан Доска</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
    <style>
        .kanban-column {
            min-height: 500px;
            background-color: #f7fafc;
            border-radius: 0.5rem;
            padding: 1rem;
        }
        .task-card {
            transition: all 0.3s ease;
            margin-bottom: 1rem;
        }
        .task-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .status-todo { border-left: 4px solid #f56565; }
        .status-in_progress { border-left: 4px solid #4299e1; }
        .status-done { border-left: 4px solid #48bb78; }
    </style>
</head>
<body class="bg-gray-100">
    <nav class="bg-white shadow-lg">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">
                <div class="flex items-center">
                    <a href="<?php echo e(route('boards.index')); ?>" class="text-xl font-bold text-gray-800">
                        <i class="fas fa-columns mr-2"></i>Канбан Доска
                    </a>
                </div>
            </div>
        </div>
    </nav>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <?php if(session('success')): ?>
            <div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
                <?php echo e(session('success')); ?>

            </div>
        <?php endif; ?>
        
        <?php echo $__env->yieldContent('content'); ?>
    </main>
</body>
</html><?php /**PATH C:\Users\danka\Desktop\code\25.12\example\resources\views/layouts/app.blade.php ENDPATH**/ ?>