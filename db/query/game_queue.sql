-- name: CreateQueue :exec 
INSERT INTO game_queue(
    user_id, status,grade,created_at,subject,lobby_id
) VALUES (
    ?,?,?,?,?,?
);

-- name: GetQueue :one
SELECT * FROM game_queue WHERE status=? AND subject=? AND grade=? AND user_id != ?  ORDER BY RAND() Limit 1;

-- name: GetOwnQueueInfo :one
SELECT * FROM game_queue WHERE user_id = ? ;

-- name: UpdateQueue :exec
UPDATE game_queue SET status=? WHERE user_id=?;

-- name: UpdateQueueLobby :exec
UPDATE game_queue SET lobby_id=? WHERE user_id=?;

-- name: RemoveQueue :exec
DELETE FROM game_queue WHERE user_id=?;