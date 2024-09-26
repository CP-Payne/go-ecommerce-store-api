-- name: GetActiveCart :one
SELECT id, user_id, status, created_at
FROM carts
WHERE user_id=$1 AND status='active'
ORDER BY created_at DESC
LIMIT 1;

-- name: GetCartItems :many
SELECT product_id, quantity
FROM cart_items
WHERE cart_id=$1;


-- name: CreateCart :exec
INSERT INTO carts (id, user_id, status)
VALUES ($1, $2, 'active');

-- name: AddItemToCart :exec
INSERT INTO cart_items (id, cart_id, product_id, quantity)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cart_id, product_id)
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity;


-- name: RemoveItemFromCart :exec
DELETE FROM cart_items
WHERE cart_id=$1 AND product_id=$2;

-- name: GetCartWithItems :many
SELECT c.id AS cart_id, c.user_id, ci.product_id, ci.quantity, p.name, p.price
FROM carts c
JOIN cart_items ci ON c.id = ci.cart_id
JOIN products p ON ci.product_id = p.id
WHERE c.id = $1;

-- name: DeleteCart :exec
DELETE FROM carts
WHERE id=$1;


-- name: ReduceItemFromCart :exec
WITH updated AS (
    UPDATE cart_items
    SET quantity = cart_items.quantity - $3
    WHERE cart_items.cart_id = $1 AND cart_items.product_id = $2
    RETURNING cart_items.quantity
)
DELETE FROM cart_items
WHERE cart_items.cart_id = $1 AND cart_items.product_id = $2 AND EXISTS (
    SELECT 1 FROM updated WHERE updated.quantity <= 0
);


