-- name: CreateChildNotification :exec 
INSERT INTO child_notifications(
   user_id, title, type, description,created_at,secondary_user_id
) VALUES (
    ?,?,?,?,?,?
);

-- name: GetChildNotification :many
SELECT * FROM child_notifications WHERE user_id=? ORDER BY created_at Desc;

-- name: DeleteChildNotification :exec
DELETE FROM child_notifications WHERE id=?;