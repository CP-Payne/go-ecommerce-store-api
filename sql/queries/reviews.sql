-- name: InsertReview :one
INSERT INTO reviews (
    id, title, review_text, rating, product_id, user_id, deleted,anonymous, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: HasUserReviewedProduct :one
SELECT EXISTS (
    SELECT 1 FROM reviews WHERE user_id = $1 AND product_id = $2 AND deleted IS NOT true
);

-- name: IsReviewOwner :one
SELECT EXISTS (
    SELECT 1 FROM reviews WHERE id = $1 AND user_id = $2 
);

-- SELECT * FROM reviews
-- WHERE product_id = $1 AND deleted IS NOT true;

-- SELECT * FROM reviews
-- WHERE user_id = $1 AND product_id = $2 AND deleted IS NOT true;

-- name: SetReviewStatusDeleted :exec
UPDATE reviews
    SET deleted = true
    WHERE user_id=$1 AND product_id=$2;

-- name: UpdateUserReview :one
UPDATE reviews
    SET title = $1, review_text = $2, rating = $3,anonymous = $4, updated_at=$5
    WHERE user_id=$6 AND product_id=$7 AND deleted IS NOT true
    RETURNING *;


-- name: GetProductReviews :many
SELECT
    r.id AS review_id,
    r.title,
    r.review_text,
    r.rating,
    r.product_id,
    CASE
        WHEN r.anonymous THEN 'anonymous'
        ELSE u.name
    END AS user_name,
    r.deleted,
    r.anonymous,
    r.created_at,
    r.updated_at
FROM
    reviews r
JOIN
    users u ON r.user_id = u.id
WHERE
    r.product_id = $1
AND deleted IS NOT true;

-- name: GetReviewByUserAndProduct :one 
SELECT
    r.id AS review_id,
    r.title,
    r.review_text,
    r.rating,
    r.product_id,
    CASE
        WHEN r.anonymous THEN 'anonymous'
        ELSE u.name
    END AS user_name,
    r.deleted,
    r.anonymous,
    r.created_at,
    r.updated_at
FROM
    reviews r
JOIN
    users u ON r.user_id = u.id
WHERE
    r.user_id = $1 AND r.product_id= $2
AND deleted IS NOT true;
