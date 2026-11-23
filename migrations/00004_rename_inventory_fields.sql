-- +goose Up
-- Preserve existing data while aligning column names with the Go and JSON models.
-- Rebuild the unique index to accommodate longer short descriptions.
ALTER TABLE `items` DROP INDEX `name_description_UNIQUE`;
ALTER TABLE `items`
  CHANGE COLUMN `value` `replacement_value` INT NOT NULL,
  CHANGE COLUMN `description` `short_description` VARCHAR(255) NOT NULL;
ALTER TABLE `items` ADD UNIQUE INDEX `name_short_description_UNIQUE` (`name`, `short_description`);

-- +goose Down
-- Reversal requires every short_description to fit within 60 characters.
ALTER TABLE `items` DROP INDEX `name_short_description_UNIQUE`;
ALTER TABLE `items`
  CHANGE COLUMN `replacement_value` `value` INT NOT NULL,
  CHANGE COLUMN `short_description` `description` VARCHAR(60) NOT NULL;
ALTER TABLE `items` ADD UNIQUE INDEX `name_description_UNIQUE` (`name`, `description`);
