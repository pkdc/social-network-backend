-- name: GetPostMembers :many
SELECT * FROM post_member
WHERE post_id = ?;

-- name: CreatePostMember :one
INSERT INTO post_member (
  user_id, post_id
) VALUES (
  ?, ?
)
RETURNING *;

-- name: DeleteMember :exec
DELETE FROM post_member
WHERE user_id = ? AND post_id = ?;