-- +goose Up
ALTER TABLE `items` ADD COLUMN `resale_value` int NOT NULL DEFAULT 0 AFTER `value`;

-- +goose Down
ALTER TABLE `items` DROP COLUMN `resale_value`;
