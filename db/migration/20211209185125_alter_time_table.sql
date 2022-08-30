-- +goose Up
-- +goose StatementBegin
ALTER TABLE time_table MODIFY start_time VARCHAR(225);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE time_table MODIFY end_time VARCHAR(225);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
