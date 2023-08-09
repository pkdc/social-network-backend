-- name: GetGroupChatNoti :many
SELECT * FROM group_chat_item
WHERE user_id = ?
ORDER BY last_msg_at DESC;

-- name: GetOneGroupChatItemByUserId :one
SELECT * FROM group_chat_item
WHERE group_id = ? AND user_id = ?;

-- name: CreateGroupChatItem :one
INSERT INTO group_chat_item (
  group_id, user_id, last_msg_at, chat_noti
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateGroupChatItem :one
UPDATE group_chat_item
SET chat_noti = ?,
last_msg_at = ?
WHERE group_id = ? AND user_id = ?
RETURNING *;

-- name: DeleteOneGroupChatItem :exec
DELETE FROM group_chat_item
WHERE group_id = ? AND user_id;

-- name: DeleteGroupChatItem :exec
DELETE FROM group_chat_item
WHERE group_id = ?;
