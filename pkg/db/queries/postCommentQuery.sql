-- name: GetPostComments :many
SELECT * FROM post_comment
WHERE post_id = ?
ORDER BY created_at;

-- name: GetAllComments :many
SELECT * FROM post_comment
ORDER BY created_at DESC;

-- name: CreatePostComment :one
INSERT INTO post_comment (
  user_id, post_id, created_at, message_, image_
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeletePostComment :exec
DELETE FROM post_comment
WHERE user_id = ? AND post_id = ?;