-- +goose Up
ALTER TABLE items ADD COLUMN info_complete BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE items DROP COLUMN info_complete;
