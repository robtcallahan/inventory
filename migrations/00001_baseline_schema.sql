-- +goose Up
-- Baseline schema for inventory: items.

CREATE TABLE IF NOT EXISTS `items` (
    `id` int NOT NULL AUTO_INCREMENT,
    `quantity` int NOT NULL,
    `name` varchar(30) NOT NULL,
    `value` int NOT NULL,
    `description` varchar(60) NOT NULL,
    `manufacturer` varchar(40) NOT NULL,
    `model` varchar(60) NOT NULL,
    `category` varchar(20) NOT NULL,
    `purchase_date` datetime DEFAULT NULL,
    `serial_number` varchar(20) DEFAULT NULL,
    `details` TEXT NOT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE INDEX `name_model_UNIQUE` (`name`, `description`)
    ) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS `items`;
