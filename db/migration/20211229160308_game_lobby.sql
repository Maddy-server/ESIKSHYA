-- +goose Up
-- +goose StatementBegin
CREATE TABLE `game_friend_lobby` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int UNIQUE NOT NULL,
  `op_id` int UNIQUE NOT NULL,
  `status` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL
  );
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `game_random_lobby` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `class` int NOT NULL,
  `op_id` int NOT NULL,
  `status` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE game_friend_lobby;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE game_random_lobby;
-- +goose StatementEnd

