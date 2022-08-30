-- name: CreateParent :exec
INSERT INTO parents (
     phone
) VALUES (
    ?
);

-- name: SaveOTP :exec
UPDATE parents SET otp=?,isVerified=0  where phone = ?;

-- name: Savepassword :exec
UPDATE parents SET password=? WHERE phone=?;

-- name: GetParent :one
SELECT id, phone, otp, isVerified,created_at FROM parents WHERE phone=?;

-- name: GetParentByRandomKey :one
SELECT parents.phone,parents.id,parents.isVerified,
    parents_detail.full_name FROM  parents_detail LEFT JOIN parents ON parents_detail.parent_id= parents.id 
    WHERE parents_detail.random_key =?;

-- name: GetParentInfo :one
SELECT * FROM parents WHERE phone=?;

-- name: SetVerification :exec
UPDATE parents SET isVerified=1 WHERE phone=?;

-- name: GetParentForLogin :one
SELECT parents.phone, parents.password,parents.id,parents.created_at,
    parents_detail.full_name, parents_detail.address FROM parents LEFT JOIN parents_detail ON parents.id = parents_detail.parent_id
    WHERE parents.phone =?;


-- name: RemoveOTP :exec
UPDATE parents SET otp="" WHERE phone=?;
