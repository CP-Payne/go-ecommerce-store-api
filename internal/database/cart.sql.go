// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: cart.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addItemToCart = `-- name: AddItemToCart :exec
INSERT INTO cart_items (id, cart_id, product_id, quantity)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cart_id, product_id)
DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
`

type AddItemToCartParams struct {
	ID        uuid.UUID
	CartID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func (q *Queries) AddItemToCart(ctx context.Context, arg AddItemToCartParams) error {
	_, err := q.db.ExecContext(ctx, addItemToCart,
		arg.ID,
		arg.CartID,
		arg.ProductID,
		arg.Quantity,
	)
	return err
}

const createCart = `-- name: CreateCart :exec
INSERT INTO carts (id, user_id, status)
VALUES ($1, $2, 'active')
`

type CreateCartParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) CreateCart(ctx context.Context, arg CreateCartParams) error {
	_, err := q.db.ExecContext(ctx, createCart, arg.ID, arg.UserID)
	return err
}

const deleteCart = `-- name: DeleteCart :exec
DELETE FROM carts
WHERE id=$1
`

func (q *Queries) DeleteCart(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteCart, id)
	return err
}

const getActiveCart = `-- name: GetActiveCart :one
SELECT id, user_id, status, created_at
FROM carts
WHERE user_id=$1 AND status='active'
ORDER BY created_at DESC
LIMIT 1
`

type GetActiveCartRow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Status    string
	CreatedAt time.Time
}

func (q *Queries) GetActiveCart(ctx context.Context, userID uuid.UUID) (GetActiveCartRow, error) {
	row := q.db.QueryRowContext(ctx, getActiveCart, userID)
	var i GetActiveCartRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getCartItems = `-- name: GetCartItems :many
SELECT product_id, quantity
FROM cart_items
WHERE cart_id=$1
`

type GetCartItemsRow struct {
	ProductID uuid.UUID
	Quantity  int32
}

func (q *Queries) GetCartItems(ctx context.Context, cartID uuid.UUID) ([]GetCartItemsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCartItems, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCartItemsRow
	for rows.Next() {
		var i GetCartItemsRow
		if err := rows.Scan(&i.ProductID, &i.Quantity); err != nil {
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

const getCartWithItems = `-- name: GetCartWithItems :many
SELECT c.id AS cart_id, c.user_id, ci.product_id, ci.quantity, p.name, p.price
FROM carts c
JOIN cart_items ci ON c.id = ci.cart_id
JOIN products p ON ci.product_id = p.id
WHERE c.id = $1
`

type GetCartWithItemsRow struct {
	CartID    uuid.UUID
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
	Name      string
	Price     string
}

func (q *Queries) GetCartWithItems(ctx context.Context, id uuid.UUID) ([]GetCartWithItemsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCartWithItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCartWithItemsRow
	for rows.Next() {
		var i GetCartWithItemsRow
		if err := rows.Scan(
			&i.CartID,
			&i.UserID,
			&i.ProductID,
			&i.Quantity,
			&i.Name,
			&i.Price,
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

const reduceItemFromCart = `-- name: ReduceItemFromCart :exec
WITH updated AS (
    UPDATE cart_items
    SET quantity = cart_items.quantity - $3
    WHERE cart_items.cart_id = $1 AND cart_items.product_id = $2
    RETURNING cart_items.quantity
)
DELETE FROM cart_items
WHERE cart_items.cart_id = $1 AND cart_items.product_id = $2 AND EXISTS (
    SELECT 1 FROM updated WHERE updated.quantity <= 0
)
`

type ReduceItemFromCartParams struct {
	CartID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int32
}

func (q *Queries) ReduceItemFromCart(ctx context.Context, arg ReduceItemFromCartParams) error {
	_, err := q.db.ExecContext(ctx, reduceItemFromCart, arg.CartID, arg.ProductID, arg.Quantity)
	return err
}

const removeItemFromCart = `-- name: RemoveItemFromCart :exec
DELETE FROM cart_items
WHERE cart_id=$1 AND product_id=$2
`

type RemoveItemFromCartParams struct {
	CartID    uuid.UUID
	ProductID uuid.UUID
}

func (q *Queries) RemoveItemFromCart(ctx context.Context, arg RemoveItemFromCartParams) error {
	_, err := q.db.ExecContext(ctx, removeItemFromCart, arg.CartID, arg.ProductID)
	return err
}
