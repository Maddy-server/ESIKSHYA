-- +goose Up
-- +goose StatementBegin
CREATE TABLE `game_queue` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int unique NOT NULL ,
  `status` varchar(255) NOT NULL,
  `subject` varchar(255) NOT NULL,
  `grade` int NOT NULL,
  `lobby_id` int NOT NULL,
  `created_at` timestamp NOT NULL,
   FOREIGN KEY(`lobby_id`) 
        REFERENCES game_random_lobby(`id`)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE game_queue;
-- +goose StatementEnd
