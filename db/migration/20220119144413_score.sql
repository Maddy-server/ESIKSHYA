-- +goose Up
-- +goose StatementBegin
CREATE TABLE `score` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `own_points` int NOT NULL,
  `op_id` int NOT NULL,
  `op_points` int NOT NULL,
  `subject` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE score;
-- +goose StatementEnd
