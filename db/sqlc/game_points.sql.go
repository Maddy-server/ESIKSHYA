// Code generated by sqlc. DO NOT EDIT.
// source: game_points.sql

package db

import (
	"context"
	"time"
)

const createScore = `-- name: CreateScore :exec
INSERT INTO game_points(
   player1_id, player2_id, player1_point,player2_point,played_time,indicator
) VALUES (
    ?,?,?,?,?,?
)
`

type CreateScoreParams struct {
	Player1ID    int32     `json:"player1_id"`
	Player2ID    int32     `json:"player2_id"`
	Player1Point int32     `json:"player1_point"`
	Player2Point int32     `json:"player2_point"`
	PlayedTime   time.Time `json:"played_time"`
	Indicator    int32     `json:"indicator"`
}

func (q *Queries) CreateScore(ctx context.Context, arg CreateScoreParams) error {
	_, err := q.exec(ctx, q.createScoreStmt, createScore,
		arg.Player1ID,
		arg.Player2ID,
		arg.Player1Point,
		arg.Player2Point,
		arg.PlayedTime,
		arg.Indicator,
	)
	return err
}

const getScore = `-- name: GetScore :one
SELECT id, player1_id, player2_id, player1_point, player2_point, indicator, played_time, deleted_at FROM game_points WHERE player1_id=? AND player2_id=? AND indicator=?
`

type GetScoreParams struct {
	Player1ID int32 `json:"player1_id"`
	Player2ID int32 `json:"player2_id"`
	Indicator int32 `json:"indicator"`
}

func (q *Queries) GetScore(ctx context.Context, arg GetScoreParams) (GamePoint, error) {
	row := q.queryRow(ctx, q.getScoreStmt, getScore, arg.Player1ID, arg.Player2ID, arg.Indicator)
	var i GamePoint
	err := row.Scan(
		&i.ID,
		&i.Player1ID,
		&i.Player2ID,
		&i.Player1Point,
		&i.Player2Point,
		&i.Indicator,
		&i.PlayedTime,
		&i.DeletedAt,
	)
	return i, err
}

const getScoreList = `-- name: GetScoreList :many
SELECT id, player1_id, player2_id, player1_point, player2_point, indicator, played_time, deleted_at FROM game_points WHERE player1_id=? AND player2_id=?
`

type GetScoreListParams struct {
	Player1ID int32 `json:"player1_id"`
	Player2ID int32 `json:"player2_id"`
}

func (q *Queries) GetScoreList(ctx context.Context, arg GetScoreListParams) ([]GamePoint, error) {
	rows, err := q.query(ctx, q.getScoreListStmt, getScoreList, arg.Player1ID, arg.Player2ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GamePoint
	for rows.Next() {
		var i GamePoint
		if err := rows.Scan(
			&i.ID,
			&i.Player1ID,
			&i.Player2ID,
			&i.Player1Point,
			&i.Player2Point,
			&i.Indicator,
			&i.PlayedTime,
			&i.DeletedAt,
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

const updateScorePlayerOne = `-- name: UpdateScorePlayerOne :exec
UPDATE game_points SET player1_point=? WHERE id=?
`

type UpdateScorePlayerOneParams struct {
	Player1Point int32 `json:"player1_point"`
	ID           int32 `json:"id"`
}

func (q *Queries) UpdateScorePlayerOne(ctx context.Context, arg UpdateScorePlayerOneParams) error {
	_, err := q.exec(ctx, q.updateScorePlayerOneStmt, updateScorePlayerOne, arg.Player1Point, arg.ID)
	return err
}

const updateScorePlayerTwo = `-- name: UpdateScorePlayerTwo :exec
UPDATE game_points SET player2_point=? WHERE id=?
`

type UpdateScorePlayerTwoParams struct {
	Player2Point int32 `json:"player2_point"`
	ID           int32 `json:"id"`
}

func (q *Queries) UpdateScorePlayerTwo(ctx context.Context, arg UpdateScorePlayerTwoParams) error {
	_, err := q.exec(ctx, q.updateScorePlayerTwoStmt, updateScorePlayerTwo, arg.Player2Point, arg.ID)
	return err
}
