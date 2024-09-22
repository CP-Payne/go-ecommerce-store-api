-- name: CreateOrder :one
INSERT INTO orders(
    id, user_id, product_total, status, payment_method, shipping_price, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;


-- name: CreateOrderItem :exec
INSERT INTO order_items(
    id, order_id, product_id, quantity, price
) VALUES ( $1, $2, $3, $4, $5);

-- name: GetOrderItemsByOrderID :many
SELECT oi.quantity, p.price, p.name, oi.product_id
FROM order_items oi
JOIN products p ON oi.product_id = p.id
WHERE oi.order_id = $1;

