-- name: GetPosts :many
SELECT * FROM post
WHERE author = ?;

-- name: GetAllPosts :many
SELECT * FROM post
ORDER BY created_at DESC;

-- name: CreatePost :one
INSERT INTO post (
  author, message_, image_, created_at, privacy
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeletePost :exec
DELETE FROM post
WHERE author = ? AND id = ?;