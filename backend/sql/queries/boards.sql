-- name: CreateBoard :one
INSERT INTO boards (title, description, owner_id)
VALUES ($1, $2, $3)
RETURNING id, title, description, owner_id, created_at, updated_at;

-- name: GetBoardByID :one
SELECT id, title, description, owner_id, created_at, updated_at
FROM boards
WHERE id = $1;

-- name: GetBoardsByOwner :many
SELECT id, title, description, owner_id, created_at, updated_at
FROM boards
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: UpdateBoard :one
UPDATE boards
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, title, description, owner_id, created_at, updated_at;

-- name: DeleteBoard :exec
DELETE FROM boards
WHERE id = $1;
