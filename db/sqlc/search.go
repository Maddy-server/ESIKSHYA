package db

import (
	"context"
)

const searchList = `
SELECT children.username, children.id FROM children LEFT JOIN children_detail 
ON children.id = children_detail.children_id
WHERE children_detail.grade=? AND 
#  Search
(
	? IS NULL OR
 (LOWER(username) LIKE ? ) ) 
`

type SearchListParams struct {
	Grade    int32  `json:"grade"`
	Username string `json:"username"`
}

type SearchListRow struct {
	Username string `json:"username"`
	ID       int32  `json:"id"`
}

func (store *SQLStore) SearchList(ctx context.Context, arg SearchListParams) ([]SearchListRow, error) {
	rows, err := store.db.Query(searchList, arg.Grade, arg.Username, arg.Username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchListRow
	for rows.Next() {
		var i SearchListRow
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
