// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: reviews.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const getProductReviews = `-- name: GetProductReviews :many
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
AND deleted IS NOT true
`

type GetProductReviewsRow struct {
	ReviewID   uuid.UUID
	Title      sql.NullString
	ReviewText sql.NullString
	Rating     int32
	ProductID  uuid.UUID
	UserName   interface{}
	Deleted    bool
	Anonymous  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (q *Queries) GetProductReviews(ctx context.Context, productID uuid.UUID) ([]GetProductReviewsRow, error) {
	rows, err := q.db.QueryContext(ctx, getProductReviews, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductReviewsRow
	for rows.Next() {
		var i GetProductReviewsRow
		if err := rows.Scan(
			&i.ReviewID,
			&i.Title,
			&i.ReviewText,
			&i.Rating,
			&i.ProductID,
			&i.UserName,
			&i.Deleted,
			&i.Anonymous,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReviewByUserAndProduct = `-- name: GetReviewByUserAndProduct :one
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
AND deleted IS NOT true
`

type GetReviewByUserAndProductParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

type GetReviewByUserAndProductRow struct {
	ReviewID   uuid.UUID
	Title      sql.NullString
	ReviewText sql.NullString
	Rating     int32
	ProductID  uuid.UUID
	UserName   interface{}
	Deleted    bool
	Anonymous  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (q *Queries) GetReviewByUserAndProduct(ctx context.Context, arg GetReviewByUserAndProductParams) (GetReviewByUserAndProductRow, error) {
	row := q.db.QueryRowContext(ctx, getReviewByUserAndProduct, arg.UserID, arg.ProductID)
	var i GetReviewByUserAndProductRow
	err := row.Scan(
		&i.ReviewID,
		&i.Title,
		&i.ReviewText,
		&i.Rating,
		&i.ProductID,
		&i.UserName,
		&i.Deleted,
		&i.Anonymous,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const hasUserReviewedProduct = `-- name: HasUserReviewedProduct :one
SELECT EXISTS (
    SELECT 1 FROM reviews WHERE user_id = $1 AND product_id = $2 AND deleted IS NOT true
)
`

type HasUserReviewedProductParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

func (q *Queries) HasUserReviewedProduct(ctx context.Context, arg HasUserReviewedProductParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, hasUserReviewedProduct, arg.UserID, arg.ProductID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const insertReview = `-- name: InsertReview :one
INSERT INTO reviews (
    id, title, review_text, rating, product_id, user_id, deleted,anonymous, created_at, updated_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, title, review_text, rating, product_id, user_id, deleted, created_at, updated_at, anonymous
`

type InsertReviewParams struct {
	ID         uuid.UUID
	Title      sql.NullString
	ReviewText sql.NullString
	Rating     int32
	ProductID  uuid.UUID
	UserID     uuid.UUID
	Deleted    bool
	Anonymous  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (q *Queries) InsertReview(ctx context.Context, arg InsertReviewParams) (Review, error) {
	row := q.db.QueryRowContext(ctx, insertReview,
		arg.ID,
		arg.Title,
		arg.ReviewText,
		arg.Rating,
		arg.ProductID,
		arg.UserID,
		arg.Deleted,
		arg.Anonymous,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.ReviewText,
		&i.Rating,
		&i.ProductID,
		&i.UserID,
		&i.Deleted,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Anonymous,
	)
	return i, err
}

const isReviewOwner = `-- name: IsReviewOwner :one
SELECT EXISTS (
    SELECT 1 FROM reviews WHERE id = $1 AND user_id = $2 
)
`

type IsReviewOwnerParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) IsReviewOwner(ctx context.Context, arg IsReviewOwnerParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, isReviewOwner, arg.ID, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const setReviewStatusDeleted = `-- name: SetReviewStatusDeleted :exec


UPDATE reviews
    SET deleted = true
    WHERE user_id=$1 AND product_id=$2
`

type SetReviewStatusDeletedParams struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}

// SELECT * FROM reviews
// WHERE product_id = $1 AND deleted IS NOT true;
// SELECT * FROM reviews
// WHERE user_id = $1 AND product_id = $2 AND deleted IS NOT true;
func (q *Queries) SetReviewStatusDeleted(ctx context.Context, arg SetReviewStatusDeletedParams) error {
	_, err := q.db.ExecContext(ctx, setReviewStatusDeleted, arg.UserID, arg.ProductID)
	return err
}

const updateUserReview = `-- name: UpdateUserReview :one
UPDATE reviews
    SET title = $1, review_text = $2, rating = $3,anonymous = $4, updated_at=$5
    WHERE user_id=$6 AND product_id=$7 AND deleted IS NOT true
    RETURNING id, title, review_text, rating, product_id, user_id, deleted, created_at, updated_at, anonymous
`

type UpdateUserReviewParams struct {
	Title      sql.NullString
	ReviewText sql.NullString
	Rating     int32
	Anonymous  bool
	UpdatedAt  time.Time
	UserID     uuid.UUID
	ProductID  uuid.UUID
}

func (q *Queries) UpdateUserReview(ctx context.Context, arg UpdateUserReviewParams) (Review, error) {
	row := q.db.QueryRowContext(ctx, updateUserReview,
		arg.Title,
		arg.ReviewText,
		arg.Rating,
		arg.Anonymous,
		arg.UpdatedAt,
		arg.UserID,
		arg.ProductID,
	)
	var i Review
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.ReviewText,
		&i.Rating,
		&i.ProductID,
		&i.UserID,
		&i.Deleted,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Anonymous,
	)
	return i, err
}
