-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    hashed_password TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
