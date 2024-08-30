-- +goose Up
CREATE TABLE customer (
    id UUID PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    hashed_password TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE customer;
