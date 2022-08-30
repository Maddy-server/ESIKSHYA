-- name: RemoveChildToken :exec
DELETE FROM child_token WHERE user_id=?;

-- name: RemoveParentsToken :exec
DELETE FROM parents_token WHERE user_id=?;