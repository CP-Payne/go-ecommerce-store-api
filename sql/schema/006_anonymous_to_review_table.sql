-- +goose Up
ALTER TABLE reviews
ADD anonymous BOOLEAN NOT NULL
DEFAULT true;

-- +goose Down
ALTER TABLE reviews
DROP COLUMN anonymous;
