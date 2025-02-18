-- +goose Up
ALTER TABLE users
ADD COLUMN is_verified BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE users
DROP COLUMN is_verified;
