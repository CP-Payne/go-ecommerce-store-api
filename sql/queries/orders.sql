-- name: CreateOrder :one
INSERT INTO orders(
    id, user_id, total_amount, status, payment_method, shipping_price, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;


-- name: CreateOrderItem :exec
INSERT INTO order_items(
    id, order_id, product_id, quantity, price
) VALUES ( $1, $2, $3, $4, $5);


