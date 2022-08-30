// Code generated by sqlc. DO NOT EDIT.
// source: general_video.sql

package db

import (
	"context"
	"database/sql"
)

const createGeneralVideo = `-- name: CreateGeneralVideo :exec
INSERT INTO general_video(
  topic,url,created_at
) VALUES (
    ?,?,?
)
`

type CreateGeneralVideoParams struct {
	Topic     sql.NullString `json:"topic"`
	Url       string         `json:"url"`
	CreatedAt sql.NullTime   `json:"created_at"`
}

func (q *Queries) CreateGeneralVideo(ctx context.Context, arg CreateGeneralVideoParams) error {
	_, err := q.exec(ctx, q.createGeneralVideoStmt, createGeneralVideo, arg.Topic, arg.Url, arg.CreatedAt)
	return err
}

const getGeneralVideo = `-- name: GetGeneralVideo :one
SELECT id, topic, url, created_at from general_video WHERE id=?
`

func (q *Queries) GetGeneralVideo(ctx context.Context, id int32) (GeneralVideo, error) {
	row := q.queryRow(ctx, q.getGeneralVideoStmt, getGeneralVideo, id)
	var i GeneralVideo
	err := row.Scan(
		&i.ID,
		&i.Topic,
		&i.Url,
		&i.CreatedAt,
	)
	return i, err
}

const getListGeneralVideo = `-- name: GetListGeneralVideo :many
SELECT id, topic, url, created_at FROM general_video ORDER BY created_at
`

func (q *Queries) GetListGeneralVideo(ctx context.Context) ([]GeneralVideo, error) {
	rows, err := q.query(ctx, q.getListGeneralVideoStmt, getListGeneralVideo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GeneralVideo
	for rows.Next() {
		var i GeneralVideo
		if err := rows.Scan(
			&i.ID,
			&i.Topic,
			&i.Url,
			&i.CreatedAt,
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