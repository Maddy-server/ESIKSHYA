-- name: CreateParentsNotification :exec 
INSERT INTO parents_notifications(
   user_id, title, type, description,created_at,secondary_user_id
) VALUES (
    ?,?,?,?,?,?
);

-- name: GetParentsNotification :many
SELECT * FROM parents_notifications WHERE user_id=? ORDER BY created_at Desc;
