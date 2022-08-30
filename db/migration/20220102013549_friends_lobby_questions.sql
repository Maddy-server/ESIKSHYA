-- +goose Up
-- +goose StatementBegin
CREATE TABLE `friends_lobby_questions` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `lobby_id` int NOT NULL,
  `questions` varchar(255) NOT NULL,
  `options_a` varchar(255) NOT NULL,
  `options_b` varchar(255) NOT NULL,
  `options_c` varchar(255) NOT NULL,
  `options_d` varchar(255) NOT NULL,
  `correct_options` varchar(255) NOT NULL,
  FOREIGN KEY(`lobby_id`) 
        REFERENCES game_friend_lobby(`id`)
        ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `random_lobby_questions` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `lobby_id` int NOT NULL,
  `questions` varchar(255) NOT NULL,
  `options_a` varchar(255) NOT NULL,
  `options_b` varchar(255) NOT NULL,
  `options_c` varchar(255) NOT NULL,
  `options_d` varchar(255) NOT NULL,
  `correct_options` varchar(255) NOT NULL,
  FOREIGN KEY(`lobby_id`) 
        REFERENCES game_random_lobby(`id`)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE lobby_questions;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE random_lobby_questions;
-- +goose StatementEnd
