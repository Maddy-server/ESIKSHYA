-- name: CreateChild :exec 
INSERT INTO children(
    parent_id, username, password 
) VALUES (
    ?,?,?
);

-- name: GetChild :one
SELECT children.id,children.isVerified, children.username, children.password, children_detail.full_name,children.created_at,
children_detail.date_of_birth, children_detail.gender, children_detail.school,
children_detail.grade,
 children_detail.country, children_detail.state FROM children LEFT JOIN children_detail
 ON children.id  = children_detail.children_id WHERE 
children.username=?;

-- name: GetChildForVerify :one
SELECT * FROM children WHERE id=?; 

-- name: SetVerificationChild :exec
UPDATE children SET isVerified=1 WHERE parent_id=?;

-- name: CheckUsernameAvailability :one
SELECT COUNT(*) FROM children WHERE username=?;

-- name: CheckChildDetail :one
SELECT COUNT(*) FROM children_detail WHERE children_id=?;

-- name: GetParentId :one
SELECT parent_id FROM children WHERE id=?;
