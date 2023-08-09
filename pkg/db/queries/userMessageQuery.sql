-- name: GetMessages :many
SELECT * FROM user_message
WHERE target_id = ? AND source_id = ? OR source_id = ? AND target_id = ?
ORDER BY created_at;

-- name: CreateMessage :one
INSERT INTO user_message (
  source_id, target_id, message_, created_at
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM user_message
WHERE source_id = ? AND target_id = ?;