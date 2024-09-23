-- name: CreateOrder :one
INSERT INTO orders(
    id, user_id, product_total,order_total, status, payment_method, shipping_price, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;


-- name: CreateOrderItem :exec
INSERT INTO order_items(
    id, order_id, product_id, quantity, price
) VALUES ( $1, $2, $3, $4, $5);

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrderByProcessorOrderID :one
SELECT * FROM orders
WHERE processor_order_id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT oi.quantity, p.price, p.name, oi.product_id
FROM order_items oi
JOIN products p ON oi.product_id = p.id
WHERE oi.order_id = $1;

-- name: SetProcessorIDAndStatus :exec
UPDATE orders
    SET processor_order_id = $1, status=$2, updated_at=$3
    WHERE id=$4;

-- name: SetOrderCompleted :exec
UPDATE orders
    SET  status=$1, payment_email=$2, payer_id=$3, updated_at=$4
    WHERE processor_order_id=$5;
