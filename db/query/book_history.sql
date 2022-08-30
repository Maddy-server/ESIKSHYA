-- name: CreateBookHistory :exec 
INSERT INTO book_history(
    book_id, user_id,created_at
) VALUES (
    ?,?,?
);

-- name: FetchBookHistoryList :many
SELECT * FROM book_history WHERE user_id=? and id>=? Limit 15; 

-- name: FetchBookHistory :one
SELECT * FROM book_history WHERE user_id=? AND book_id=?;

-- name: UpdateBookHistory :exec
UPDATE book_history SET created_at=? WHERE id=?;