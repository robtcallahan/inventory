-- +goose Up
ALTER TABLE items
  ADD COLUMN room VARCHAR(100) NOT NULL DEFAULT '',
  ADD COLUMN `condition` VARCHAR(30) NOT NULL DEFAULT '',
  ADD COLUMN purchase_source VARCHAR(100) NOT NULL DEFAULT '',
  ADD COLUMN purchase_price INT NOT NULL DEFAULT 0,
  ADD COLUMN website_url VARCHAR(2048) NOT NULL DEFAULT '',
  ADD COLUMN manual_url VARCHAR(2048) NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE items
  DROP COLUMN manual_url,
  DROP COLUMN website_url,
  DROP COLUMN purchase_price,
  DROP COLUMN purchase_source,
  DROP COLUMN `condition`,
  DROP COLUMN room;
