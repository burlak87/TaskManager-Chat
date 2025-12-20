-- name: CreateTwoFACode :one
INSERT INTO two_fa_codes (user_id, code, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetActiveTwoFACode :one
SELECT * FROM two_fa_codes 
WHERE user_id = $1 
  AND is_used = false 
  AND attempts < 3 
  AND expires_at > NOW()
ORDER BY created_at DESC 
LIMIT 1;

-- name: UpdateTwoFAAttempts :exec
UPDATE two_fa_codes SET attempts = $1 WHERE id = $2;

-- name: MarkTwoFACodeUsed :exec
UPDATE two_fa_codes SET is_used = true WHERE id = $1;

-- name: GetRecentCodeRequests :one
SELECT COUNT(*) FROM two_fa_codes 
WHERE user_id = $1 AND created_at > $2;

-- name: GetRecentVerificationAttempts :one
SELECT COUNT(*) FROM two_fa_codes 
WHERE user_id = $1 
  AND created_at > $2 
  AND (attempts > 0 OR is_used = true);