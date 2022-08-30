// Code generated by sqlc. DO NOT EDIT.
// source: game_questions.sql

package db

import (
	"context"
)

const createQuestions = `-- name: CreateQuestions :exec
INSERT INTO game_questions(
   class, subject, questions,options_a,options_b,options_c,options_d,correct_options,difficulty_level
) VALUES (
    ?,?,?,?,?,?,?,?,?
)
`

type CreateQuestionsParams struct {
	Class           int32  `json:"class"`
	Subject         string `json:"subject"`
	Questions       string `json:"questions"`
	OptionsA        string `json:"options_a"`
	OptionsB        string `json:"options_b"`
	OptionsC        string `json:"options_c"`
	OptionsD        string `json:"options_d"`
	CorrectOptions  string `json:"correct_options"`
	DifficultyLevel int32  `json:"difficulty_level"`
}

func (q *Queries) CreateQuestions(ctx context.Context, arg CreateQuestionsParams) error {
	_, err := q.exec(ctx, q.createQuestionsStmt, createQuestions,
		arg.Class,
		arg.Subject,
		arg.Questions,
		arg.OptionsA,
		arg.OptionsB,
		arg.OptionsC,
		arg.OptionsD,
		arg.CorrectOptions,
		arg.DifficultyLevel,
	)
	return err
}

const getQuestions = `-- name: GetQuestions :many
SELECT id, class, subject, questions, options_a, options_b, options_c, options_d, correct_options, difficulty_level FROM game_questions WHERE class=? AND subject=?  ORDER BY RAND() Limit 10
`

type GetQuestionsParams struct {
	Class   int32  `json:"class"`
	Subject string `json:"subject"`
}

func (q *Queries) GetQuestions(ctx context.Context, arg GetQuestionsParams) ([]GameQuestion, error) {
	rows, err := q.query(ctx, q.getQuestionsStmt, getQuestions, arg.Class, arg.Subject)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GameQuestion
	for rows.Next() {
		var i GameQuestion
		if err := rows.Scan(
			&i.ID,
			&i.Class,
			&i.Subject,
			&i.Questions,
			&i.OptionsA,
			&i.OptionsB,
			&i.OptionsC,
			&i.OptionsD,
			&i.CorrectOptions,
			&i.DifficultyLevel,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
