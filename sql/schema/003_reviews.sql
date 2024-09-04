-- +goose Up
CREATE TABLE reviews(
    id UUID PRIMARY KEY,
    title VARCHAR(255),
    review_text TEXT,
    rating INTEGER NOT NULL DEFAULT 5,
    product_id UUID NOT NULL REFERENCES products(id),
    user_id UUID NOT NULL REFERENCES users(id),
    deleted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- +goose Down
DROP TABLE reviews;
