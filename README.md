# Git
``` Bash
# Клонировать репозиторий
git clone https://github.com/your-team/taskmanager-chat.git
cd taskmanager-chat

# Установить pre-commit хуки
pre-commit install

# Получить последние изменения
git pull origin main

# Создать новую ветку для задачи
git checkout -b feature/ваша-фича

# После завершения работы
git add .
git commit -m "feat: описание изменения"
git push origin feature/ваша-фича

# Создать Pull Request через GitHub UI
```

# Работаем с Github так
```
1. git checkout main
2. git pull origin main
3. git checkout -b feature/type-description
   # Типы: feat, fix, docs, style, refactor, test, chore
   # Пример: feature/auth-login
4. Работа над задачей
5. git add .
6. git commit -m "feat: добавить вход через JWT"
7. git push origin feature/auth-login
8. Создать PR на GitHub
9. Ждать ревью от 2 разработчиков
10. После аппрува - мердж в main
```

# Если нужно обновить ветку с main:
```
git checkout feature/ваша-ветка
git fetch origin
git merge origin/main
```

# Docker
``` Bash
# Запустить весь проект (фронтенд + бэкенд + БД)
docker-compose up --build

# Остановить все контейнеры
docker-compose down

# Перезапустить конкретный сервис
docker-compose restart backend

# Просмотреть логи
docker-compose logs -f backend

# Очистить всё (контейнеры, образы, тома)
docker system prune -a --volumes

# Проверить состояние контейнеров
docker-compose ps
```

# Ежедневно для всех
``` Bash
# Утром:
git pull origin main
npm install (если фронтенд)
go mod tidy (если бэкенд)

# В течение дня:
# Работа в своей ветке
# Частые коммиты

# Вечером:
git push origin ваша-ветка
# Проверить CI пайплайн на GitHub
```

# Список всех веток для команды
branches=(
  "feature/auth"          # FE-1 + BE-1
  "feature/kanban"        # FE-2 + BE-2
  "feature/chat"         # FE-3 + BE-3
  "feature/notifications" # FE-4 + BE-4
  "feature/ui-polish"    # FE-5
  "feature/initial-setup" # Для инициализации проекта
)

# Создаём и переключаемся на каждую ветку
for branch in "${branches[@]}"; do
  git checkout -b "$branch"
  git push -u origin "$branch"
done

# Возвращаемся на main
git checkout main