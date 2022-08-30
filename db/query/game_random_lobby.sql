-- name: CreateGameRandomLobby :exec 
INSERT INTO game_random_lobby(
   user_id, op_id, status,created_at,class
) VALUES (
    ?,?,?,?,?
);

-- name: GetGameRandomLobby :one
SELECT * FROM game_random_lobby WHERE user_id=? AND op_id=? ;

-- name: GetGameRandomLobbyById :one
SELECT * FROM game_random_lobby WHERE id=? ;

-- name: GetFakeGameRandomLobbyByClass :one
SELECT * FROM game_random_lobby WHERE class=? ;

-- name: UpdateGameRandomLobby :exec
UPDATE game_random_lobby SET status=? WHERE id=?;

-- name: DeleteGameRandomLobby :exec
DELETE FROM game_random_lobby WHERE id=?;


