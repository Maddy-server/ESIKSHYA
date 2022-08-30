-- name: AddTimeTable :exec
INSERT INTO time_table(
    children_id, class, section, description, day, start_time, end_time
) VALUES (
    ?,?,?,?,?,?,?
);

-- name: GetTimeTable :many
SELECT * FROM time_table WHERE children_id=?;

-- name: RemoveTimeTableAll :exec
DELETE FROM time_table WHERE children_id=?;

-- name: RemoveTimeTableByDescription :exec
DELETE FROM time_table WHERE description=?;
