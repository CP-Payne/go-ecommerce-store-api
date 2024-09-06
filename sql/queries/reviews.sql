-- name: InsertReview :one
INSERT INTO reviews (
    id, title, review_text, rating, product_id, user_id, deleted, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: HasUserReviewedProduct :one
SELECT EXISTS (
    SELECT 1 FROM reviews WHERE user_id = $1 AND product_id = $2 AND deleted IS NOT true
);

-- name: IsReviewOwner :one
SELECT EXISTS (
    SELECT 1 FROM reviews WHERE id = $1 AND user_id = $2 
);

-- name: GetProductReviews :many
SELECT * FROM reviews
WHERE product_id = $1 AND deleted IS NOT true;

-- name: GetReviewByUserAndProduct :one 
SELECT * FROM reviews
WHERE user_id = $1 AND product_id = $2 AND deleted IS NOT true;

-- name: SetReviewStatusDeleted :exec
UPDATE reviews
    SET deleted = true
    WHERE user_id=$1 AND product_id=$2;

-- name: UpdateUserReview :one
UPDATE reviews
    SET title = $1, review_text = $2, rating = $3, updated_at=$4
    WHERE user_id=$5 AND product_id=$6 AND deleted IS NOT true
    RETURNING *;
