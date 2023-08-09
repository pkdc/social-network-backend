-- name: GetFollowers :many
SELECT * FROM user_follower
WHERE target_id = ?;

-- name: GetFollowings :many
SELECT * FROM user_follower
WHERE source_id = ?;

-- name: CheckFollower :one
SELECT * FROM user_follower
WHERE source_id = ? AND target_id = ? AND (status_= 1 OR status_ = 2);

-- name: CreateFollower :one
INSERT INTO user_follower (
  source_id, target_id, status_
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: DeleteFollower :exec
DELETE FROM user_follower
WHERE source_id = ? AND target_id = ?;

-- name: ReplyFollowReq :exec
UPDATE user_follower
set status_ = 1
WHERE source_id = ? AND target_id = ?;

-- name: UpdateFollower :one
UPDATE user_follower
set status_ = ?
WHERE source_id = ? AND target_id = ?
RETURNING *;