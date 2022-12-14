// Code generated by sqlc. DO NOT EDIT.
// source: book_count.sql

package db

import (
	"context"
	"time"
)

const addCount = `-- name: AddCount :exec
INSERT INTO book_count(
    book_id, count,created_at
) VALUES (
    ?,?,?
)
`

type AddCountParams struct {
	BookID    int32     `json:"book_id"`
	Count     int32     `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

func (q *Queries) AddCount(ctx context.Context, arg AddCountParams) error {
	_, err := q.exec(ctx, q.addCountStmt, addCount, arg.BookID, arg.Count, arg.CreatedAt)
	return err
}

const fetchBookCount = `-- name: FetchBookCount :one
SELECT count FROM book_count WHERE book_id=?
`

func (q *Queries) FetchBookCount(ctx context.Context, bookID int32) (int32, error) {
	row := q.queryRow(ctx, q.fetchBookCountStmt, fetchBookCount, bookID)
	var count int32
	err := row.Scan(&count)
	return count, err
}

const updateBookCount = `-- name: UpdateBookCount :exec
UPDATE book_count SET count=? WHERE book_id=?
`

type UpdateBookCountParams struct {
	Count  int32 `json:"count"`
	BookID int32 `json:"book_id"`
}

func (q *Queries) UpdateBookCount(ctx context.Context, arg UpdateBookCountParams) error {
	_, err := q.exec(ctx, q.updateBookCountStmt, updateBookCount, arg.Count, arg.BookID)
	return err
}
