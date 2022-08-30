-- name: CreateQuestions :exec 
INSERT INTO game_questions(
   class, subject, questions,options_a,options_b,options_c,options_d,correct_options,difficulty_level
) VALUES (
    ?,?,?,?,?,?,?,?,?
);

-- name: GetQuestions :many
SELECT * FROM game_questions WHERE class=? AND subject=?  ORDER BY RAND() Limit 10 ;
