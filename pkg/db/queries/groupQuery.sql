-- name: GetGroup :one
SELECT * FROM group_
WHERE id = ? LIMIT 1;

-- name: GetAllGroups :many
SELECT * FROM group_;

-- name: CreateGroup :one
INSERT INTO group_ (
  title, creator, description_, created_at
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM group_
WHERE id = ?;

-- name: CheckIfCreator :one
SELECT COUNT(*) FROM group_
WHERE creator = ? AND id = ? LIMIT 1;