-- +goose Up
-- +goose StatementBegin
ALTER TABLE video ADD COLUMN img_url VARCHAR(225);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE video ADD COLUMN video_id VARCHAR(225);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
