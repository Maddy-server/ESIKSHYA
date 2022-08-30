-- +goose Up
-- +goose StatementBegin
CREATE TABLE `general_video` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `topic` varchar(255),
  `url` varchar(255) NOT NULL,
  `created_at` timestamp DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE general_video;
-- +goose StatementEnd
