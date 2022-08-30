// Code generated by sqlc. DO NOT EDIT.
// source: friends.sql

package db

import (
	"context"
	"database/sql"
)

const acceptFriendRequest = `-- name: AcceptFriendRequest :exec
UPDATE friends SET status=? WHERE id=?
`

type AcceptFriendRequestParams struct {
	Status string `json:"status"`
	ID     int32  `json:"id"`
}

func (q *Queries) AcceptFriendRequest(ctx context.Context, arg AcceptFriendRequestParams) error {
	_, err := q.exec(ctx, q.acceptFriendRequestStmt, acceptFriendRequest, arg.Status, arg.ID)
	return err
}

const checkFriendsList = `-- name: CheckFriendsList :many
SELECT children.username, children.id FROM children LEFT JOIN friends 
ON children.id=friends.sender_id OR children.id=friends.receiver_id 
WHERE children.id!=? AND (friends.receiver_id=? OR friends.sender_id=?)
`

type CheckFriendsListParams struct {
	ID         int32 `json:"id"`
	ReceiverID int32 `json:"receiver_id"`
	SenderID   int32 `json:"sender_id"`
}

type CheckFriendsListRow struct {
	Username string `json:"username"`
	ID       int32  `json:"id"`
}

func (q *Queries) CheckFriendsList(ctx context.Context, arg CheckFriendsListParams) ([]CheckFriendsListRow, error) {
	rows, err := q.query(ctx, q.checkFriendsListStmt, checkFriendsList, arg.ID, arg.ReceiverID, arg.SenderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CheckFriendsListRow
	for rows.Next() {
		var i CheckFriendsListRow
		if err := rows.Scan(&i.Username, &i.ID); err != nil {
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

const getFriend = `-- name: GetFriend :one
SELECT id, sender_id, receiver_id, status, friends_at FROM friends WHERE sender_id=? AND receiver_id=? AND status=?
`

type GetFriendParams struct {
	SenderID   int32  `json:"sender_id"`
	ReceiverID int32  `json:"receiver_id"`
	Status     string `json:"status"`
}

func (q *Queries) GetFriend(ctx context.Context, arg GetFriendParams) (Friend, error) {
	row := q.queryRow(ctx, q.getFriendStmt, getFriend, arg.SenderID, arg.ReceiverID, arg.Status)
	var i Friend
	err := row.Scan(
		&i.ID,
		&i.SenderID,
		&i.ReceiverID,
		&i.Status,
		&i.FriendsAt,
	)
	return i, err
}

const getFriendsList = `-- name: GetFriendsList :many
SELECT children.username, children.id FROM children LEFT JOIN friends 
ON children.id=friends.sender_id OR children.id=friends.receiver_id 
WHERE friends.status=? And children.id!=? AND (friends.receiver_id=? OR friends.sender_id=?)
`

type GetFriendsListParams struct {
	Status     string `json:"status"`
	ID         int32  `json:"id"`
	ReceiverID int32  `json:"receiver_id"`
	SenderID   int32  `json:"sender_id"`
}

type GetFriendsListRow struct {
	Username string `json:"username"`
	ID       int32  `json:"id"`
}

func (q *Queries) GetFriendsList(ctx context.Context, arg GetFriendsListParams) ([]GetFriendsListRow, error) {
	rows, err := q.query(ctx, q.getFriendsListStmt, getFriendsList,
		arg.Status,
		arg.ID,
		arg.ReceiverID,
		arg.SenderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFriendsListRow
	for rows.Next() {
		var i GetFriendsListRow
		if err := rows.Scan(&i.Username, &i.ID); err != nil {
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

const rejectFriendRequest = `-- name: RejectFriendRequest :exec
DELETE FROM friends WHERE id=?
`

func (q *Queries) RejectFriendRequest(ctx context.Context, id int32) error {
	_, err := q.exec(ctx, q.rejectFriendRequestStmt, rejectFriendRequest, id)
	return err
}

const sendFriendRequest = `-- name: SendFriendRequest :exec
INSERT INTO friends(
    sender_id, receiver_id, status, friends_at
) VALUES (
    ?,?,?,?
)
`

type SendFriendRequestParams struct {
	SenderID   int32        `json:"sender_id"`
	ReceiverID int32        `json:"receiver_id"`
	Status     string       `json:"status"`
	FriendsAt  sql.NullTime `json:"friends_at"`
}

func (q *Queries) SendFriendRequest(ctx context.Context, arg SendFriendRequestParams) error {
	_, err := q.exec(ctx, q.sendFriendRequestStmt, sendFriendRequest,
		arg.SenderID,
		arg.ReceiverID,
		arg.Status,
		arg.FriendsAt,
	)
	return err
}
