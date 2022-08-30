-- name: AddCount :exec 
INSERT INTO book_count(
    book_id, count,created_at
) VALUES (
    ?,?,?
);

-- name: UpdateBookCount :exec
UPDATE book_count SET count=? WHERE book_id=?;

-- name: FetchBookCount :one
SELECT count FROM book_count WHERE book_id=? ;

