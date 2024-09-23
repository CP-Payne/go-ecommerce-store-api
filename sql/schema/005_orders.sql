-- +goose Up
CREATE TABLE orders(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    processor_order_id VARCHAR(100),
    product_total DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'created',
    order_total DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(255) NOT NULL,
    payment_email VARCHAR(255),
    payer_id VARCHAR(255), 
    shipping_price DECIMAL(10, 2) NOT NULL,
    cart_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price DECIMAL(10, 2) NOT NULL,
    UNIQUE (order_id, product_id)
);


-- +goose Down
DROP TABLE order_items;
DROP TABLE orders;
