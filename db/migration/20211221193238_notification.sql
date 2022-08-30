-- +goose Up
-- +goose StatementBegin
CREATE TABLE `child_notifications` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `secondary_user_id` int ,
  `title` varchar(255) NOT NULL, 
  `type` varchar(255) NOT NULL, 
  `description` varchar(255) NOT NULL, 
  `created_at` timestamp NOT NULL
);
-- +goose StatementBegin
-- +goose StatementEnd
CREATE TABLE `parents_notifications` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `secondary_user_id` int ,
  `title` varchar(255) NOT NULL, 
  `type` varchar(255) NOT NULL, 
  `description` varchar(255) NOT NULL, 
  `created_at` timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE child_notifications;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE parents_notifications;
-- +goose StatementEnd