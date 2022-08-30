-- name: CreateParentDetail :exec
INSERT INTO parents_detail(
    parent_id, full_name, address,random_key
) VALUES (
    ?,?,?,?
);

-- name: CompairKey :one 
 SELECT * from parents_detail WHERE random_key=?;

-- name: GetParentDetail :one
SELECT * from parents_detail WHERE parent_id=?;


-- name: EditParentDetail :exec
UPDATE parents_detail SET full_name=?, address=? WHERE parent_id=?;

