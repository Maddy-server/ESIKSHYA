-- name: CreateGameFriendLobby :exec 
INSERT INTO game_friend_lobby(
   user_id, op_id, status,created_at
) VALUES (
    ?,?,?,?
);

-- name: GetGameFriendLobby :one
SELECT * FROM game_friend_lobby WHERE user_id=? AND op_id=? ;

-- name: UpdateGameFriendLobby :exec
UPDATE game_friend_lobby SET status=? WHERE id=?;

-- name: DeleteGameFriendLobby :exec
DELETE FROM game_friend_lobby WHERE id=?;