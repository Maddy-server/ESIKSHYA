// Code generated by sqlc. DO NOT EDIT.
// source: book_history.sql

package db

import (
	"context"
	"time"
)

const createBookHistory = `-- name: CreateBookHistory :exec
INSERT INTO book_history(
    book_id, user_id,created_at
) VALUES (
    ?,?,?
)
`

type CreateBookHistoryParams struct {
	BookID    int32     `json:"book_id"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) CreateBookHistory(ctx context.Context, arg CreateBookHistoryParams) error {
	_, err := q.exec(ctx, q.createBookHistoryStmt, createBookHistory, arg.BookID, arg.UserID, arg.CreatedAt)
	return err
}

const fetchBookHistory = `-- name: FetchBookHistory :one
SELECT id, book_id, user_id, created_at FROM book_history WHERE user_id=? AND book_id=?
`

type FetchBookHistoryParams struct {
	UserID int32 `json:"user_id"`
	BookID int32 `json:"book_id"`
}

func (q *Queries) FetchBookHistory(ctx context.Context, arg FetchBookHistoryParams) (BookHistory, error) {
	row := q.queryRow(ctx, q.fetchBookHistoryStmt, fetchBookHistory, arg.UserID, arg.BookID)
	var i BookHistory
	err := row.Scan(
		&i.ID,
		&i.BookID,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const fetchBookHistoryList = `-- name: FetchBookHistoryList :many
SELECT id, book_id, user_id, created_at FROM book_history WHERE user_id=? and id>=? Limit 15
`

type FetchBookHistoryListParams struct {
	UserID int32 `json:"user_id"`
	ID     int32 `json:"id"`
}

func (q *Queries) FetchBookHistoryList(ctx context.Context, arg FetchBookHistoryListParams) ([]BookHistory, error) {
	rows, err := q.query(ctx, q.fetchBookHistoryListStmt, fetchBookHistoryList, arg.UserID, arg.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BookHistory
	for rows.Next() {
		var i BookHistory
		if err := rows.Scan(
			&i.ID,
			&i.BookID,
			&i.UserID,
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

const updateBookHistory = `-- name: UpdateBookHistory :exec
UPDATE book_history SET created_at=? WHERE id=?
`

type UpdateBookHistoryParams struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int32     `json:"id"`
}

func (q *Queries) UpdateBookHistory(ctx context.Context, arg UpdateBookHistoryParams) error {
	_, err := q.exec(ctx, q.updateBookHistoryStmt, updateBookHistory, arg.CreatedAt, arg.ID)
	return err
}
