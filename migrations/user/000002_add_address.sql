-- +goose Up
ALTER TABLE users ADD COLUMN address TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE users DROP COLUMN address;
