-- name: CreatePrivateChatItem :one
 INSERT INTO private_chat_item (
  source_id, target_id, chat_noti, last_msg_at
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: GetPrivateChatItem :many
SELECT * FROM private_chat_item
WHERE target_id = ?
ORDER BY last_msg_at DESC;

-- name: GetOnePrivateChatItem :one
SELECT * FROM private_chat_item
WHERE source_id = ? AND target_id = ?;

-- name: DeletePrivateChatItem :exec
DELETE FROM private_chat_item
WHERE source_id = ? AND target_id = ?;

-- name: UpdatePrivateChatItem :one
UPDATE private_chat_item
SET chat_noti = ?,
last_msg_at = ?
WHERE source_id = ? AND target_id = ?
RETURNING *;