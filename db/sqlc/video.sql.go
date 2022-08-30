// Code generated by sqlc. DO NOT EDIT.
// source: video.sql

package db

import (
	"context"
	"database/sql"
)

const createVideo = `-- name: CreateVideo :exec
INSERT INTO video(
  class,subject,topic,url,created_at,img_url,video_id
) VALUES (
    ?,?,?,?,?,?,?
)
`

type CreateVideoParams struct {
	Class     int32          `json:"class"`
	Subject   string         `json:"subject"`
	Topic     sql.NullString `json:"topic"`
	Url       string         `json:"url"`
	CreatedAt sql.NullTime   `json:"created_at"`
	ImgUrl    sql.NullString `json:"img_url"`
	VideoID   sql.NullString `json:"video_id"`
}

func (q *Queries) CreateVideo(ctx context.Context, arg CreateVideoParams) error {
	_, err := q.exec(ctx, q.createVideoStmt, createVideo,
		arg.Class,
		arg.Subject,
		arg.Topic,
		arg.Url,
		arg.CreatedAt,
		arg.ImgUrl,
		arg.VideoID,
	)
	return err
}

const getClassVideo = `-- name: GetClassVideo :many
SELECT id, class, subject, topic, url, created_at, img_url, video_id FROM video WHERE class=? ORDER BY created_at ASC
`

func (q *Queries) GetClassVideo(ctx context.Context, class int32) ([]Video, error) {
	rows, err := q.query(ctx, q.getClassVideoStmt, getClassVideo, class)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Video
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.ID,
			&i.Class,
			&i.Subject,
			&i.Topic,
			&i.Url,
			&i.CreatedAt,
			&i.ImgUrl,
			&i.VideoID,
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

const getClassVideoFree = `-- name: GetClassVideoFree :many
SELECT id, class, subject, topic, url, created_at, img_url, video_id FROM video WHERE class=? ORDER BY created_at ASC LIMIT 2
`

func (q *Queries) GetClassVideoFree(ctx context.Context, class int32) ([]Video, error) {
	rows, err := q.query(ctx, q.getClassVideoFreeStmt, getClassVideoFree, class)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Video
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.ID,
			&i.Class,
			&i.Subject,
			&i.Topic,
			&i.Url,
			&i.CreatedAt,
			&i.ImgUrl,
			&i.VideoID,
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

const getSubjectVideo = `-- name: GetSubjectVideo :many
SELECT id, class, subject, topic, url, created_at, img_url, video_id FROM video WHERE class=? AND subject=? ORDER BY created_at ASC
`

type GetSubjectVideoParams struct {
	Class   int32  `json:"class"`
	Subject string `json:"subject"`
}

func (q *Queries) GetSubjectVideo(ctx context.Context, arg GetSubjectVideoParams) ([]Video, error) {
	rows, err := q.query(ctx, q.getSubjectVideoStmt, getSubjectVideo, arg.Class, arg.Subject)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Video
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.ID,
			&i.Class,
			&i.Subject,
			&i.Topic,
			&i.Url,
			&i.CreatedAt,
			&i.ImgUrl,
			&i.VideoID,
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

const getSubjectVideoFree = `-- name: GetSubjectVideoFree :many
SELECT id, class, subject, topic, url, created_at, img_url, video_id FROM video WHERE class=? AND subject=? ORDER BY created_at ASC LIMIT 2
`

type GetSubjectVideoFreeParams struct {
	Class   int32  `json:"class"`
	Subject string `json:"subject"`
}

func (q *Queries) GetSubjectVideoFree(ctx context.Context, arg GetSubjectVideoFreeParams) ([]Video, error) {
	rows, err := q.query(ctx, q.getSubjectVideoFreeStmt, getSubjectVideoFree, arg.Class, arg.Subject)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Video
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.ID,
			&i.Class,
			&i.Subject,
			&i.Topic,
			&i.Url,
			&i.CreatedAt,
			&i.ImgUrl,
			&i.VideoID,
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

const getVideo = `-- name: GetVideo :one
SELECT id, class, subject, topic, url, created_at, img_url, video_id from video WHERE id=?
`

func (q *Queries) GetVideo(ctx context.Context, id int32) (Video, error) {
	row := q.queryRow(ctx, q.getVideoStmt, getVideo, id)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Class,
		&i.Subject,
		&i.Topic,
		&i.Url,
		&i.CreatedAt,
		&i.ImgUrl,
		&i.VideoID,
	)
	return i, err
}
