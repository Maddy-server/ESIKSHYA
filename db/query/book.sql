-- name: AddBook :exec 
INSERT INTO book(
    book_name, writer, section, randomunique,created_at
) VALUES (
    ?,?,?,?,?
);

-- name: FetchBookAfterCreated :one
SELECT * FROM book WHERE randomunique=?;

-- name: UpdateBook :exec
UPDATE book SET content=? WHERE id=?;

-- name: FetchBookListHome :many
SELECT book.id,book.book_name, book_count.count FROM book JOIN book_count ON book.id = book_count.book_id WHERE book.id>=? Limit 20;

-- name: FetchBookListBySection :many
SELECT book.id,book.book_name, book_count.count FROM book JOIN book_count ON book.id = book_count.book_id WHERE book.section=? And book.id>=? Limit 20;

-- name: FetchPopularBookListBySection :many
SELECT book.id, book.book_name, book_count.count FROM book JOIN book_count ON book.id = book_count.book_id WHERE book.section=? ORDER BY book_count.count DESC Limit 5;

-- name: FetchBookById :one
SELECT book.id, book.book_name, book_count.count FROM book JOIN book_count ON book.id = book_count.book_id WHERE book.id=? ;

-- name: FetchPopularBook :many
SELECT book.id, book.book_name, book_count.count FROM book JOIN book_count ON book.id = book_count.book_id  ORDER BY book_count.count DESC Limit 5;

-- name: FetchNewBook :many
SELECT book.id, book.book_name,book_count.count FROM book JOIN book_count ON book.id = book_count.book_id ORDER BY book.created_at DESC Limit 5 ;

-- name: FetchBookDetailsById :one
SELECT book.id, book.book_name,book.writer,book.description,book_count.count  FROM book JOIN book_count ON book.id = book_count.book_id WHERE book.id=? ;

-- name: FetchBookContent :one
SELECT content FROM book WHERE id=? ;

