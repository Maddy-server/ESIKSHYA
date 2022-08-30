-- +goose Up
-- +goose StatementBegin
CREATE TABLE `book` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `book_name` varchar(255) NOT NULL,
  `content` longtext ,
  `writer` varchar(255) NOT NULL,
  `section` varchar(255) NOT NULL,
  `randomunique` varchar(255) unique NOT NULL,
  `description` text NOT NULL,
  `created_at` timestamp NOT NULL,
  `deleted_at` timestamp DEFAULT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `book_count` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `book_id` int unique NOT NULL,
  `count` int NOT NULL,
  `created_at` timestamp NOT NULL,
  `deleted_at` timestamp DEFAULT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `book_rating` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `book_id` int NOT NULL,
  `rating` int NOT NULL,
  `user_id` int NOT NULL,
  `created_at` timestamp NOT NULL,
  `deleted_at` timestamp DEFAULT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `book_saved` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `book_id` int unique NOT NULL,
  `user_id` int NOT NULL,
  `created_at` timestamp NOT NULL
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE `book_history` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `book_id` int unique NOT NULL ,
  `user_id` int NOT NULL,
  `created_at` timestamp NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE book;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE book_count;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE book_rating;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE book_saved;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE book_history;
-- +goose StatementEnd

