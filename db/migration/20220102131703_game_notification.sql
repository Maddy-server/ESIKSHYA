-- +goose Up
-- +goose StatementBegin
CREATE TABLE `game_notifications` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `oponent_id` int NOT NULL,
  `title` varchar(255) NOT NULL, 
  `type` varchar(255) NOT NULL, 
  `description` varchar(255) NOT NULL, 
  `subject` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  `grade` int NOT NULL,
  `created_at` timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE game_notifications;
-- +goose StatementEnd
