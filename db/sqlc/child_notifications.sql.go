// Code generated by sqlc. DO NOT EDIT.
// source: child_notifications.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createChildNotification = `-- name: CreateChildNotification :exec
INSERT INTO child_notifications(
   user_id, title, type, description,created_at,secondary_user_id
) VALUES (
    ?,?,?,?,?,?
)
`

type CreateChildNotificationParams struct {
	UserID          int32         `json:"user_id"`
	Title           string        `json:"title"`
	Type            string        `json:"type"`
	Description     string        `json:"description"`
	CreatedAt       time.Time     `json:"created_at"`
	SecondaryUserID sql.NullInt32 `json:"secondary_user_id"`
}

func (q *Queries) CreateChildNotification(ctx context.Context, arg CreateChildNotificationParams) error {
	_, err := q.exec(ctx, q.createChildNotificationStmt, createChildNotification,
		arg.UserID,
		arg.Title,
		arg.Type,
		arg.Description,
		arg.CreatedAt,
		arg.SecondaryUserID,
	)
	return err
}

const deleteChildNotification = `-- name: DeleteChildNotification :exec
DELETE FROM child_notifications WHERE id=?
`

func (q *Queries) DeleteChildNotification(ctx context.Context, id int32) error {
	_, err := q.exec(ctx, q.deleteChildNotificationStmt, deleteChildNotification, id)
	return err
}

const getChildNotification = `-- name: GetChildNotification :many
SELECT id, user_id, secondary_user_id, title, type, description, created_at FROM child_notifications WHERE user_id=? ORDER BY created_at Desc
`

func (q *Queries) GetChildNotification(ctx context.Context, userID int32) ([]ChildNotification, error) {
	rows, err := q.query(ctx, q.getChildNotificationStmt, getChildNotification, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ChildNotification
	for rows.Next() {
		var i ChildNotification
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SecondaryUserID,
			&i.Title,
			&i.Type,
			&i.Description,
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