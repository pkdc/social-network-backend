-- name: GetGroupPosts :many
SELECT * FROM group_post
WHERE group_id = ?
ORDER BY created_at;

-- name: GetGroupPostById :one
SELECT * FROM group_post
WHERE id = ?;

-- name: CreateGroupPost :one
INSERT INTO group_post (
  author, group_id, message_, image_, created_at
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteGroupPost :exec
DELETE FROM group_post
WHERE group_id = ? AND author = ? AND id = ?;