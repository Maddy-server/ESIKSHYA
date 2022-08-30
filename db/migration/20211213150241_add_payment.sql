-- +goose Up
-- +goose StatementBegin
CREATE TABLE `payment` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `transaction_id` varchar(255) NOT NULL, 
  `transaction_token` varchar(255) NOT NULL, 
  `method` varchar(255) NOT NULL, 
  `parent_id` int NOT NULL,
  `child_id` int NOT NULL,
  `amount` int NOT NULL, 
  `pay_at` timestamp DEFAULT now(),
  `expire_at` timestamp
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE `payment_number` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `number` varchar(255) NOT NULL, 
  `method` varchar(255) NOT NULL,
  `save` tinyint(1) DEFAULT 1, 
  `parent_id` int NOT NULL
);
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
DROP TABLE payment;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE payment_number;
-- +goose StatementEnd
