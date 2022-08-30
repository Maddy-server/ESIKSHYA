-- name: CreateScore :exec 
INSERT INTO game_points(
   player1_id, player2_id, player1_point,player2_point,played_time,indicator
) VALUES (
    ?,?,?,?,?,?
);

-- name: GetScore :one 
SELECT * FROM game_points WHERE player1_id=? AND player2_id=? AND indicator=?;

-- name: GetScoreList :many 
SELECT * FROM game_points WHERE player1_id=? AND player2_id=?;

-- name: UpdateScorePlayerOne :exec
UPDATE game_points SET player1_point=? WHERE id=?;

-- name: UpdateScorePlayerTwo :exec
UPDATE game_points SET player2_point=? WHERE id=?;



