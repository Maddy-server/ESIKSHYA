-- name: SendFriendRequest :exec 
INSERT INTO friends(
    sender_id, receiver_id, status, friends_at
) VALUES (
    ?,?,?,?
);

-- name: AcceptFriendRequest :exec
UPDATE friends SET status=? WHERE id=?;

-- name: GetFriend :one
SELECT * FROM friends WHERE sender_id=? AND receiver_id=? AND status=?;

-- name: GetFriendsList :many
SELECT children.username, children.id FROM children LEFT JOIN friends 
ON children.id=friends.sender_id OR children.id=friends.receiver_id 
WHERE friends.status=? And children.id!=? AND (friends.receiver_id=? OR friends.sender_id=?);

-- name: CheckFriendsList :many
SELECT children.username, children.id FROM children LEFT JOIN friends 
ON children.id=friends.sender_id OR children.id=friends.receiver_id 
WHERE children.id!=? AND (friends.receiver_id=? OR friends.sender_id=?);

-- name: RejectFriendRequest :exec
DELETE FROM friends WHERE id=?;