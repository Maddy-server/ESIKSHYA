// Code generated by sqlc. DO NOT EDIT.
// source: children_detail.sql

package db

import (
	"context"
	"database/sql"
)

const createChildDetail = `-- name: CreateChildDetail :exec
INSERT INTO children_detail(
    children_id, full_name, date_of_birth,grade, gender,school,country,state
) VALUES (
    ?,?,?,?,?,?,?,?
)
`

type CreateChildDetailParams struct {
	ChildrenID  int32          `json:"children_id"`
	FullName    string         `json:"full_name"`
	DateOfBirth string         `json:"date_of_birth"`
	Grade       int32          `json:"grade"`
	Gender      string         `json:"gender"`
	School      string         `json:"school"`
	Country     sql.NullString `json:"country"`
	State       sql.NullString `json:"state"`
}

func (q *Queries) CreateChildDetail(ctx context.Context, arg CreateChildDetailParams) error {
	_, err := q.exec(ctx, q.createChildDetailStmt, createChildDetail,
		arg.ChildrenID,
		arg.FullName,
		arg.DateOfBirth,
		arg.Grade,
		arg.Gender,
		arg.School,
		arg.Country,
		arg.State,
	)
	return err
}

const editChildDetail = `-- name: EditChildDetail :exec
UPDATE children_detail SET full_name=?,grade=?, gender=?,school=?,country=?,state=? WHERE children_id=?
`

type EditChildDetailParams struct {
	FullName   string         `json:"full_name"`
	Grade      int32          `json:"grade"`
	Gender     string         `json:"gender"`
	School     string         `json:"school"`
	Country    sql.NullString `json:"country"`
	State      sql.NullString `json:"state"`
	ChildrenID int32          `json:"children_id"`
}

func (q *Queries) EditChildDetail(ctx context.Context, arg EditChildDetailParams) error {
	_, err := q.exec(ctx, q.editChildDetailStmt, editChildDetail,
		arg.FullName,
		arg.Grade,
		arg.Gender,
		arg.School,
		arg.Country,
		arg.State,
		arg.ChildrenID,
	)
	return err
}

const getChildDetail = `-- name: GetChildDetail :one
SELECT id, children_id, full_name, grade, date_of_birth, gender, school, country, state from children_detail WHERE children_id=?
`

func (q *Queries) GetChildDetail(ctx context.Context, childrenID int32) (ChildrenDetail, error) {
	row := q.queryRow(ctx, q.getChildDetailStmt, getChildDetail, childrenID)
	var i ChildrenDetail
	err := row.Scan(
		&i.ID,
		&i.ChildrenID,
		&i.FullName,
		&i.Grade,
		&i.DateOfBirth,
		&i.Gender,
		&i.School,
		&i.Country,
		&i.State,
	)
	return i, err
}

const getChildDetailListOnCountry = `-- name: GetChildDetailListOnCountry :many
SELECT id, children_id, full_name, grade, date_of_birth, gender, school, country, state from children_detail WHERE country=?
`

func (q *Queries) GetChildDetailListOnCountry(ctx context.Context, country sql.NullString) ([]ChildrenDetail, error) {
	rows, err := q.query(ctx, q.getChildDetailListOnCountryStmt, getChildDetailListOnCountry, country)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ChildrenDetail
	for rows.Next() {
		var i ChildrenDetail
		if err := rows.Scan(
			&i.ID,
			&i.ChildrenID,
			&i.FullName,
			&i.Grade,
			&i.DateOfBirth,
			&i.Gender,
			&i.School,
			&i.Country,
			&i.State,
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

const getChildDetailListOnState = `-- name: GetChildDetailListOnState :many
SELECT id, children_id, full_name, grade, date_of_birth, gender, school, country, state from children_detail WHERE state=?
`

func (q *Queries) GetChildDetailListOnState(ctx context.Context, state sql.NullString) ([]ChildrenDetail, error) {
	rows, err := q.query(ctx, q.getChildDetailListOnStateStmt, getChildDetailListOnState, state)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ChildrenDetail
	for rows.Next() {
		var i ChildrenDetail
		if err := rows.Scan(
			&i.ID,
			&i.ChildrenID,
			&i.FullName,
			&i.Grade,
			&i.DateOfBirth,
			&i.Gender,
			&i.School,
			&i.Country,
			&i.State,
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

const getChildrenDetails = `-- name: GetChildrenDetails :many
SELECT children.id, children.username, children_detail.grade, children_detail.full_name,
children_detail.date_of_birth, children_detail.gender, children_detail.school,
 children_detail.country, children_detail.state FROM children LEFT JOIN children_detail
 ON children.id  = children_detail.children_id WHERE children.parent_id=?
`

type GetChildrenDetailsRow struct {
	ID          int32          `json:"id"`
	Username    string         `json:"username"`
	Grade       sql.NullInt32  `json:"grade"`
	FullName    sql.NullString `json:"full_name"`
	DateOfBirth sql.NullString `json:"date_of_birth"`
	Gender      sql.NullString `json:"gender"`
	School      sql.NullString `json:"school"`
	Country     sql.NullString `json:"country"`
	State       sql.NullString `json:"state"`
}

func (q *Queries) GetChildrenDetails(ctx context.Context, parentID int32) ([]GetChildrenDetailsRow, error) {
	rows, err := q.query(ctx, q.getChildrenDetailsStmt, getChildrenDetails, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetChildrenDetailsRow
	for rows.Next() {
		var i GetChildrenDetailsRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Grade,
			&i.FullName,
			&i.DateOfBirth,
			&i.Gender,
			&i.School,
			&i.Country,
			&i.State,
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