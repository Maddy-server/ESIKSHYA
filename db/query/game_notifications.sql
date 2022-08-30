-- name: CreateGameNotification :exec 
INSERT INTO game_notifications(
   user_id, title, type,oponent_id, description,created_at,subject,status,grade
) VALUES (
    ?,?,?,?,?,?,?,?,?
);

-- name: GetGameNotification :many
SELECT * FROM game_notifications WHERE user_id=? ORDER BY created_at Desc;

-- name: DeleteGameNotification :exec
DELETE FROM game_notifications WHERE id=? AND user_id=? ;