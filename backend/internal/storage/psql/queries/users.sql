-- name: CreateUser :one
INSERT INTO users (firstname, lastname, email, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users 
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY username 
LIMIT $1 OFFSET $2;

-- name: GetUserWithOrders :many
SELECT 
    u.*,
    o.id as order_id,
    o.total_amount,
    o.status
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE u.id = $1;

-- name: GetUserStats :one
SELECT 
    COUNT(*) as total_users,
    AVG(age) as average_age,
    MIN(created_at) as oldest_user,
    MAX(created_at) as newest_user
FROM users;

-- name: UpdateUser :one
UPDATE users 
SET username = $2, email = $3, age = $4, updated_at = NOW()
WHERE id = $1 
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;

-- name: UpdateTwoFAStatus :exec
UPDATE users SET two_fa_enabled = $1 WHERE id = $2;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE token = $1;

-- name: CreateLoginAttempt :one
INSERT INTO login_attempts (email, result, attempt_time)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRecentFailedAttempts :one
SELECT COUNT(*) FROM login_attempts 
WHERE email = $1 
  AND result = false 
  AND attempt_time >= $2;

-- name: GetBlockedStatus :one
SELECT blocked_until FROM login_attempts 
WHERE email = $1 AND blocked_until >= $2 
ORDER BY blocked_until DESC LIMIT 1;

-- name: BlockUser :exec
UPDATE login_attempts 
SET blocked_until = $2 
WHERE email = $1 AND blocked_until IS NULL;