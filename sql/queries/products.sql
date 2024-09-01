-- name: ListProducts :many
SELECT * FROM products
WHERE (created_at > $1 OR (created_at = $1 AND id > $2))
ORDER BY created_at, id
LIMIT $3;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1;

-- name: GetTotalProducts :one
SELECT COUNT(*) FROM products;

-- name: GetAllProducts :many
SELECT * FROM products;
