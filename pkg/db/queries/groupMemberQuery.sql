-- name: GetGroupMembers :many
SELECT * FROM group_member
WHERE group_id = ? AND status_ = ?;


-- name: GetGroupMembersByUserId :many
SELECT * FROM group_member
WHERE user_id = ? AND status_ = ?;

-- name: CheckIfMember :one
SELECT COUNT(*) FROM group_member
WHERE group_id = ? AND user_id = ? AND status_ = ? LIMIT 1;

-- name: GetGroupMembersByGroupId :many
SELECT user.* FROM group_member JOIN user ON group_member.user_id = user.id
WHERE group_member.group_id = ? AND group_member.status_ = ?;

-- name: GetGroupMembersByGroupIdWithoutStatus :many
SELECT user.* FROM group_member JOIN user ON group_member.user_id = user.id
WHERE group_member.group_id = ? ;

-- name: GetAllGroupsByUser :many
SELECT group_.* FROM group_member JOIN group_ ON group_member.group_id = group_.id
WHERE group_member.user_id = ? AND group_member.status_ = 1;



-- name: CreateGroupMember :one
INSERT INTO group_member (
  user_id, group_id, status_
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateGroupMember :one
UPDATE group_member
set status_ = ?
WHERE group_id = ? AND user_id = ?
RETURNING *;

-- name: DeleteGroupMember :exec
DELETE FROM group_member
WHERE group_id = ? AND user_id = ?;

-- name: UpdateGroupMemberChatNoti :one
UPDATE group_member
set chat_noti = ?
WHERE group_id = ? AND user_id = ?
RETURNING *;

-- name: UpdateGroupMemberChatNotiUnseen :many
UPDATE group_member
set chat_noti = 0
WHERE group_id = ?
RETURNING *;

-- name: UpdateGroupMemberChatNotiSeen :one
UPDATE group_member
set chat_noti = 1
WHERE group_id = ? and user_id = ?
RETURNING *;
