-- name: CreateRandomLobbyQuestions :exec 
INSERT INTO random_lobby_questions(
   lobby_id, questions,options_a,options_b,options_c,options_d,correct_options
) VALUES (
    ?,?,?,?,?,?,?
);

-- name: GetRandomLobbyQuestions :many
SELECT * FROM random_lobby_questions WHERE lobby_id = ?;

