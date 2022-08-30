-- name: CreateChildDetail :exec
INSERT INTO children_detail(
    children_id, full_name, date_of_birth,grade, gender,school,country,state
) VALUES (
    ?,?,?,?,?,?,?,?
);

-- name: GetChildDetail :one
SELECT * from children_detail WHERE children_id=?;

-- name: GetChildDetailListOnState :many
SELECT * from children_detail WHERE state=?;

-- name: GetChildDetailListOnCountry :many
SELECT * from children_detail WHERE country=?;


-- name: EditChildDetail :exec
UPDATE children_detail SET full_name=?,grade=?, gender=?,school=?,country=?,state=? WHERE children_id=?;


-- name: GetChildrenDetails :many
SELECT children.id, children.username, children_detail.grade, children_detail.full_name,
children_detail.date_of_birth, children_detail.gender, children_detail.school,
 children_detail.country, children_detail.state FROM children LEFT JOIN children_detail
 ON children.id  = children_detail.children_id WHERE children.parent_id=?;