-- name: CreateSaveBook :exec 
INSERT INTO book_saved(
    book_id, user_id,created_at
) VALUES (
    ?,?,?
);

-- name: FetchSavedBookList :many
SELECT * FROM book_saved WHERE user_id=? and id>=? Limit 15;

-- name: FetchSavedBook :one
SELECT * FROM book_saved WHERE user_id=? AND book_id=?;

-- name: RemovedSavedBook :exec
DELETE FROM book_saved WHERE id=? ;