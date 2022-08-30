-- name: CreateVideo :exec 
INSERT INTO video(
  class,subject,topic,url,created_at,img_url,video_id
) VALUES (
    ?,?,?,?,?,?,?
);

-- name: GetVideo :one
SELECT * from video WHERE id=?;

-- name: GetClassVideo :many
SELECT * FROM video WHERE class=? ORDER BY created_at ASC;

-- name: GetSubjectVideo :many
SELECT * FROM video WHERE class=? AND subject=? ORDER BY created_at ASC ;

-- name: GetClassVideoFree :many
SELECT * FROM video WHERE class=? ORDER BY created_at ASC LIMIT 2;

-- name: GetSubjectVideoFree :many
SELECT * FROM video WHERE class=? AND subject=? ORDER BY created_at ASC LIMIT 2;