-- name: CreateGeneralVideo :exec 
INSERT INTO general_video(
  topic,url,created_at
) VALUES (
    ?,?,?
);

-- name: GetGeneralVideo :one
SELECT * from general_video WHERE id=?;

-- name: GetListGeneralVideo :many
SELECT * FROM general_video ORDER BY created_at ;