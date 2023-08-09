-- name: GetGroupEvent :one
SELECT * FROM group_event
WHERE id = ? AND group_id = ?;

-- name: GetGroupEventById :one
SELECT * FROM group_event
WHERE id = ? ;

-- name: GetGroupEvents :many
SELECT * FROM group_event
WHERE group_id = ?
ORDER BY created_at;

-- name: GetGroupEventsWithoutId :many
SELECT * FROM group_event
ORDER BY created_at;

-- name: CreateGroupEvent :one
INSERT INTO group_event (
  author, group_id, title, description_, created_at, date_
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteGroupEvent :exec
DELETE FROM group_event
WHERE group_id = ? AND id = ?;