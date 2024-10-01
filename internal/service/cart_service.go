package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/models"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CartService struct {
	logger *zap.Logger
	db     *database.Queries
}

func NewCartService(db *database.Queries) *CartService {
	return &CartService{
		logger: config.GetLogger(),
		db:     db,
	}
}

// Get cart with items
func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (models.Cart, error) {

	logger := s.logger.With(
		zap.String("method", "GetCart"),
		zap.String("userID", userID.String()),
	)

	cart, err := s.getActiveCart(ctx, userID)
	if err != nil {
		// If ErrNotFound then there is no active cart for user, create new one
		if errors.Is(err, apperrors.ErrNotFound) {
			logger.Info("creating new user cart")
			cart, err = s.createCart(ctx, userID)
			if err != nil {
				logger.Error("failed to create user cart", zap.Error(err))
				return models.Cart{}, fmt.Errorf("failed to retrieve user cart: %w", err)
			}
		} else {
			logger.Error("failed to retrieve and determine if user has an active cart", zap.Error(err))
			return models.Cart{}, fmt.Errorf("failed to retrieve cart: %w", err)

		}
	}

	cartWithItems, err := s.db.GetCartWithItems(ctx, cart.ID)
	if err != nil {
		logger.Error("failed to retrieve cart items", zap.Error(err), zap.String("cartID", cart.ID.String()))
		return cart, fmt.Errorf("failed to retrieve cart items: %w", err)
	}

	itemsInfo := make([]models.CartItem, 0, len(cartWithItems))

	for _, cartItem := range cartWithItems {
		p, err := strconv.ParseFloat(cartItem.Price, 32)
		if err != nil {
			logger.Error("failed to convert cart item price (string) to (float)", zap.Error(err),
				zap.String("productID", cartItem.CartID.String()), zap.String("price", cartItem.Price))
			return cart, fmt.Errorf("failed to convert item prices to float: %w", err)
		}

		itemsInfo = append(itemsInfo, models.CartItem{
			ProductID: cartItem.ProductID,
			Quantity:  int(cartItem.Quantity),
			Price:     float32(p),
			Name:      cartItem.Name,
		})
	}

	cart.Items = itemsInfo
	return cart, nil
}

func (s *CartService) createCart(ctx context.Context, userID uuid.UUID) (models.Cart, error) {
	logger := s.logger.With(
		zap.String("method", "createCart"),
		zap.String("userID", userID.String()),
	)
	cart := models.Cart{
		ID:     uuid.New(),
		UserID: userID,
		Items:  []models.CartItem{},
	}

	err := s.db.CreateCart(ctx, database.CreateCartParams{
		ID:     cart.ID,
		UserID: cart.UserID,
	})
	if err != nil {
		logger.Error("failed to create new user cart", zap.Error(err))
		return models.Cart{}, fmt.Errorf("failed to create user cart: %w", err)
	}
	logger.Info("new user cart created", zap.String("cartID", cart.ID.String()))
	return cart, nil
}

func (s *CartService) CreateTemporaryProductCart(ctx context.Context, userID uuid.UUID, product models.Product, quantity int) models.Cart {
	cart := models.Cart{
		ID:     uuid.New(),
		UserID: userID,
		Status: "temporary",
		Items: []models.CartItem{
			{
				ProductID: product.ID,
				Quantity:  quantity,
				Price:     product.Price,
				Name:      product.Name,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return cart
}

func (s *CartService) DeleteCart(ctx context.Context, cartID uuid.UUID) error {
	logger := s.logger.With(
		zap.String("method", "DeleteCart"),
		zap.String("cartID", cartID.String()),
	)

	err := s.db.DeleteCart(ctx, cartID)
	if err != nil {
		logger.Error("failed to delete cart", zap.Error(err))
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	logger.Info("successfully deleted cart")
	return nil
}

func (s *CartService) AddToCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error {

	logger := s.logger.With(
		zap.String("method", "AddToCart"),
		zap.String("userID", userID.String()),
	)

	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		logger.Error("failed to retrieve user cart", zap.Error(err))
		return fmt.Errorf("failed to retrieve user cart: %w", err)
	}

	// Add item to cart
	err = s.db.AddItemToCart(ctx, database.AddItemToCartParams{
		ID:        uuid.New(),
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  int32(quantity),
	})
	if err != nil {
		logger.Error("failed to add item to cart", zap.Error(err), zap.String("cartID", cart.ID.String()))
		return fmt.Errorf("failed to add item to cart: %w", err)
	}

	logger.Info("successfully added item to cart",
		zap.String("cartID", cart.ID.String()),
		zap.String("productID", productID.String()))
	return nil
}

func (s *CartService) ReduceFromCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error {

	logger := s.logger.With(
		zap.String("method", "ReduceFromCart"),
		zap.String("userID", userID.String()),
		zap.String("productID", productID.String()),
		zap.Int("quantity", quantity),
	)

	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		logger.Error("failed to retrieve user cart", zap.Error(err))
		return fmt.Errorf("failed to retrieve user cart: %w", err)
	}

	// Reduce item from cart
	err = s.db.ReduceItemFromCart(ctx, database.ReduceItemFromCartParams{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  int32(quantity),
	})
	if err == nil {
		// If item was succesfully reduced with quantity still being greater than 0
		logger.Info("succesfully reduced item in cart", zap.String("cartID", cart.ID.String()))
		return nil
	}

	if apperrors.IsCheckViolation(err) {
		// Reducing quantity resulted in violation (quantity < 1), remove the produt from cart
		err = s.db.RemoveItemFromCart(ctx, database.RemoveItemFromCartParams{
			CartID:    cart.ID,
			ProductID: productID,
		})
		if err != nil {
			logger.Error("failed to remove product from cart", zap.Error(err), zap.String("cartID", cart.ID.String()))
			return fmt.Errorf("failed to remove item from cart: %w", err)
		}
		logger.Info("succesfully removed product from cart (quantity < 1)", zap.String("cartID", cart.ID.String()))
		return nil
	}
	logger.Error("failed to reduce cart item quantity", zap.Error(err), zap.String("cartID", cart.ID.String()))
	return fmt.Errorf("failed to reduce cart item quantity: %w", err)
}

func (s *CartService) RemoveFromCart(ctx context.Context, userID, productID uuid.UUID) error {
	logger := s.logger.With(
		zap.String("method", "RemoveFromCart"),
		zap.String("userID", userID.String()),
		zap.String("productID", productID.String()),
	)
	cart, err := s.getActiveCart(ctx, userID)
	if err != nil {
		logger.Error("failed to retrieve user cart", zap.Error(err))
		return fmt.Errorf("failed to retrieve user cart: %w", err)
	}
	return s.db.RemoveItemFromCart(ctx, database.RemoveItemFromCartParams{
		CartID:    cart.ID,
		ProductID: productID,
	})
}

func (s *CartService) getActiveCart(ctx context.Context, userID uuid.UUID) (models.Cart, error) {
	logger := s.logger.With(
		zap.String("method", "getActiveCart"),
		zap.String("userID", userID.String()),
	)

	// Get active cart
	cartRecord, err := s.db.GetActiveCart(ctx, userID)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			logger.Info("no active cart found for user", zap.Error(err))
			return models.Cart{}, fmt.Errorf("no active cart found for user: %w", apperrors.ErrNotFound)
		}
		logger.Error("failed to fetch user cart", zap.Error(err))
		return models.Cart{}, fmt.Errorf("failed to fetch user cart: %w", err)
	}

	cart := models.Cart{
		ID:     cartRecord.ID,
		UserID: userID,
		Status: cartRecord.Status,
	}

	// Fetch cart items
	cartItems, err := s.db.GetCartItems(ctx, cart.ID)
	if err != nil {
		logger.Error("failed to fetch cart items", zap.Error(err), zap.String("cartID", cart.ID.String()))
		return models.Cart{}, fmt.Errorf("failed to fetch cart items: %w", err)
	}

	// Map cart items to the Cart model
	for _, item := range cartItems {
		cart.Items = append(cart.Items, models.CartItem{
			ProductID: item.ProductID,
			Quantity:  int(item.Quantity),
		})
	}

	return cart, nil
}
