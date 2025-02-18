-- +goose Up
ALTER TABLE users
ADD COLUMN verification_code TEXT NOT NULL DEFAULT "unverified";

-- +goose Down
ALTER TABLE users
DROP COLUMN verification_code;