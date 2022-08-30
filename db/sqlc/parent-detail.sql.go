// Code generated by sqlc. DO NOT EDIT.
// source: parent-detail.sql

package db

import (
	"context"
)

const compairKey = `-- name: CompairKey :one
 SELECT id, random_key, parent_id, full_name, address from parents_detail WHERE random_key=?
`

func (q *Queries) CompairKey(ctx context.Context, randomKey string) (ParentsDetail, error) {
	row := q.queryRow(ctx, q.compairKeyStmt, compairKey, randomKey)
	var i ParentsDetail
	err := row.Scan(
		&i.ID,
		&i.RandomKey,
		&i.ParentID,
		&i.FullName,
		&i.Address,
	)
	return i, err
}

const createParentDetail = `-- name: CreateParentDetail :exec
INSERT INTO parents_detail(
    parent_id, full_name, address,random_key
) VALUES (
    ?,?,?,?
)
`

type CreateParentDetailParams struct {
	ParentID  int32  `json:"parent_id"`
	FullName  string `json:"full_name"`
	Address   string `json:"address"`
	RandomKey string `json:"random_key"`
}

func (q *Queries) CreateParentDetail(ctx context.Context, arg CreateParentDetailParams) error {
	_, err := q.exec(ctx, q.createParentDetailStmt, createParentDetail,
		arg.ParentID,
		arg.FullName,
		arg.Address,
		arg.RandomKey,
	)
	return err
}

const editParentDetail = `-- name: EditParentDetail :exec
UPDATE parents_detail SET full_name=?, address=? WHERE parent_id=?
`

type EditParentDetailParams struct {
	FullName string `json:"full_name"`
	Address  string `json:"address"`
	ParentID int32  `json:"parent_id"`
}

func (q *Queries) EditParentDetail(ctx context.Context, arg EditParentDetailParams) error {
	_, err := q.exec(ctx, q.editParentDetailStmt, editParentDetail, arg.FullName, arg.Address, arg.ParentID)
	return err
}

const getParentDetail = `-- name: GetParentDetail :one
SELECT id, random_key, parent_id, full_name, address from parents_detail WHERE parent_id=?
`

func (q *Queries) GetParentDetail(ctx context.Context, parentID int32) (ParentsDetail, error) {
	row := q.queryRow(ctx, q.getParentDetailStmt, getParentDetail, parentID)
	var i ParentsDetail
	err := row.Scan(
		&i.ID,
		&i.RandomKey,
		&i.ParentID,
		&i.FullName,
		&i.Address,
	)
	return i, err
}
