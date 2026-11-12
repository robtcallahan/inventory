-- +goose Up
-- Baseline schema for inventory: items.

CREATE TABLE IF NOT EXISTS `items` (
    `id` int NOT NULL AUTO_INCREMENT,
    `quantity` int NOT NULL,
    `name` varchar(100) NOT NULL,
    `manufacturer` varchar(40) NOT NULL,
    `model` varchar(120) NOT NULL,
    `description` TEXT NOT NULL,
    `serial_number` varchar(20) DEFAULT NULL,
    `category` varchar(20) NOT NULL,
    `purchase_date` datetime DEFAULT NULL,
    `value` int NOT NULL,
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name_UNIQUE` (`name`)
    ) ENGINE=InnoDB;

-- +goose Down
DROP TABLE IF EXISTS `items`;
