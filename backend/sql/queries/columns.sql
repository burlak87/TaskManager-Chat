-- name: CreateColumn :one
INSERT INTO columns (board_id, title, position)
VALUES ($1, $2, $3)
RETURNING id, board_id, title, position, created_at, updated_at;

-- name: GetColumnsByBoardID :many
SELECT id, board_id, title, position, created_at, updated_at
FROM columns
WHERE board_id = $1
ORDER BY position;

-- name: GetColumnByID :one
SELECT id, board_id, title, position, created_at, updated_at
FROM columns
WHERE id = $1;

-- name: UpdateColumn :one
UPDATE columns
SET title = COALESCE($2, title),
    position = COALESCE($3, position),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, board_id, title, position, created_at, updated_at;

-- name: DeleteColumn :exec
DELETE FROM columns
WHERE id = $1;
