-- +goose Up
-- +goose StatementBegin
CREATE TABLE `parents` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `phone` varchar(255) UNIQUE NOT NULL,
  `password` varchar(255) ,
  `otp` varchar(255),
  `isVerified` tinyint(1) DEFAULT 0,
  `created_at` timestamp DEFAULT now(),
  `deleted_at` timestamp
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `parents_detail` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `random_key` varchar(255) UNIQUE NOT NULL,
  `parent_id` int UNIQUE NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `children` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `parent_id` int NOT NULL,
  `username` varchar(255) UNIQUE NOT NULL,
  `password` varchar(255) NOT NULL,
  `isVerified` tinyint(1) DEFAULT 0,
  `created_at` timestamp DEFAULT now(),
  `deleted_at` timestamp
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `children_detail` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `children_id` int UNIQUE NOT NULL,
  `full_name` varchar(255) NOT NULL,
  `grade` int NOT NULL,
  `date_of_birth`varchar(255)  NOT NULL,
  `gender` varchar(255) NOT NULL,
  `school` varchar(255) NOT NULL,
  `country` varchar(255),
  `state` varchar(255)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `time_table` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `children_id` int NOT NULL,
  `class` int NOT NULL,
  `section` varchar(255) NOT NULL,
  `description` varchar(255) NOT NULL,
  `day` varchar(255) NOT NULL,
  `start_time` time NOT NULL,
  `end_time` time NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `game_points` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `player1_id` int NOT NULL,
  `player2_id` int NOT NULL,
  `player1_point` int NOT NULL,
  `player2_point` int NOT NULL,
  `indicator` int NOT NULL,
  `played_time` timestamp NOT NULL,
  `Deleted_At` timestamp
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `game_questions` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `class` int NOT NULL,
  `subject` varchar(255) NOT NULL,
  `questions` varchar(255) NOT NULL,
  `options_a` varchar(255) NOT NULL,
  `options_b` varchar(255) NOT NULL,
  `options_c` varchar(255) NOT NULL,
  `options_d` varchar(255) NOT NULL,
  `correct_options` varchar(255) NOT NULL,
  `difficulty_level` int NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `video` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `class` int NOT NULL,
  `subject` varchar(255) NOT NULL,
  `topic` varchar(255),
  `url` varchar(255) NOT NULL,
  `created_at` timestamp DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE `parents_detail` ADD FOREIGN KEY (`parent_id`) REFERENCES `parents` (`id`);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE `children` ADD FOREIGN KEY (`parent_id`) REFERENCES `parents` (`id`);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE `children_detail` ADD FOREIGN KEY (`children_id`) REFERENCES `children` (`id`);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE `time_table` ADD FOREIGN KEY (`children_id`) REFERENCES `children` (`id`);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE game_points;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE game_questions;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE video;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE time_table;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE children_detail;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE children;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE parents_detail;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE parents;
-- +goose StatementEnd
