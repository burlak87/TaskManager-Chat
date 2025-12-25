INSERT INTO notifications (
    user_id,
    task_id,
    title,
    message,
    type,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC;

UPDATE notifications
SET is_read = TRUE
WHERE id = $1 AND user_id = $2;

SELECT COUNT(*) FROM notifications
WHERE user_id = $1 AND is_read = FALSE;
