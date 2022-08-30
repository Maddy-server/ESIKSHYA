-- +goose Up
-- +goose StatementBegin
CREATE TABLE `child_token` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int UNIQUE NOT NULL,
  `token` varchar(255) NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `parents_token` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int UNIQUE NOT NULL,
  `token` varchar(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE child_token;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE parents_token;
-- +goose StatementEnd
