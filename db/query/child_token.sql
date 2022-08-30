-- name: CreateChildToken :exec 
INSERT INTO child_token(
    user_id, token
) VALUES (
    ?,?
);

-- name: UpdateChildToken :exec
UPDATE child_token SET token=? WHERE user_id=?;

-- name: GetChildToken :one
SELECT * FROM child_token WHERE user_id=?;

