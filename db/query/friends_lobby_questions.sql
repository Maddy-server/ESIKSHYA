-- name: CreateFriendsLobbyQuestions :exec 
INSERT INTO friends_lobby_questions(
   lobby_id, questions,options_a,options_b,options_c,options_d,correct_options
) VALUES (
    ?,?,?,?,?,?,?
);

-- name: GetFriendsLobbyQuestions :many
SELECT * FROM friends_lobby_questions WHERE lobby_id = ?;

