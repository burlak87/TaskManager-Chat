-- name: CreateTask :one
INSERT INTO tasks (board_id, column_id, title, description, position)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, board_id, column_id, title, description, position, created_at, updated_at;

-- name: GetTasksByBoardID :many
SELECT id, board_id, column_id, title, description, position, created_at, updated_at
FROM tasks
WHERE board_id = $1
ORDER BY position;

-- name: GetTaskByID :one
SELECT id, board_id, column_id, title, description, position, created_at, updated_at
FROM tasks
WHERE id = $1;

-- name: UpdateTask :one
UPDATE tasks
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    column_id = COALESCE($4, column_id),
    position = COALESCE($5, position),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, board_id, column_id, title, description, position, created_at, updated_at;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
