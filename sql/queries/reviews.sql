-- name: InsertReview :one
INSERT INTO reviews (
    id, title, review_text, rating, product_id, user_id, deleted, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;
