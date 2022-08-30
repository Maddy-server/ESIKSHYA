-- name: CreateParentsToken :exec 
INSERT INTO parents_token(
    user_id, token
) VALUES (
    ?,?
);

-- name: UpdateParentsToken :exec
UPDATE parents_token SET token=? WHERE user_id=?;

-- name: GetParentsToken :one
SELECT * FROM parents_token WHERE user_id=?;
