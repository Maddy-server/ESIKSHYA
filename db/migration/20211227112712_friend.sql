-- +goose Up
-- +goose StatementBegin
CREATE TABLE `friends` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `sender_id` int NOT NULL,
  `receiver_id` int NOT NULL,
  `status` varchar(255) NOT NULL,
  `friends_at` timestamp DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE friends;
-- +goose StatementEnd
