-- name: GetGroupEventMembers :many
SELECT * FROM group_event_member
WHERE event_id = ?;

-- name: GetGroupEventMember :one
SELECT COUNT(*) FROM group_event_member
WHERE event_id = ? AND user_id = ? LIMIT 1;

-- name: GetGroupEventMembersGoing :many
SELECT * FROM group_event_member
WHERE event_id = ? AND status_ = 1;

-- name: CreateGroupEventMember :one
INSERT INTO group_event_member (
  user_id, event_id, status_
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateGroupEventMember :one
UPDATE group_event_member
set status_ = ?
WHERE event_id = ? AND user_id = ?
RETURNING *;


-- name: ExecUpdateGroupEventMember :exec
UPDATE group_event_member
SET status_ = CASE
    WHEN status_ = 0 THEN 1
    ELSE status_
END
WHERE event_id = ? AND user_id = ?;

-- name: GetGroupEventsByUserNoReply :many
SELECT * FROM group_event_member
WHERE user_id = ? AND status_ = 0;


-- name: GetGroupEventsByUserAccepted :many
SELECT * FROM group_event_member
WHERE user_id = ? AND status_ = 2;

-- name: DeleteGroupEventMember :exec
DELETE FROM group_event_member
WHERE event_id = ? AND user_id = ?;