package service

import (
	"context"
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

// TODO: Add better logging and errors
// TODO: Add ability to reduce quantity and add quantity to a product in cart

// Get cart with items
func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (models.Cart, error) {
	cart, err := s.getActiveCart(ctx, userID)
	// If err != nil then there is no cart, create one
	if err != nil {
		cart = models.Cart{
			ID:     uuid.New(),
			UserID: userID,
			Items:  []models.CartItem{},
		}

		err = s.db.CreateCart(ctx, database.CreateCartParams{
			ID:     cart.ID,
			UserID: cart.UserID,
		})
		if err != nil {
			return models.Cart{}, err
		}
	}

	cartWithItems, err := s.db.GetCartWithItems(ctx, cart.ID)
	if err != nil {
		return cart, err
	}

	itemsInfo := make([]models.CartItem, 0, len(cartWithItems))

	for _, cartItem := range cartWithItems {

		p, err := strconv.ParseFloat(cartItem.Price, 32)
		if err != nil {
			s.logger.Error("failed to convert string price to float", zap.Error(err))
			return cart, err
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
	err := s.db.DeleteCart(ctx, cartID)
	if err != nil {
		s.logger.Error("failed to delete cart", zap.Error(err), zap.String("cartID", cartID.String()))
		return fmt.Errorf("failed to delete cart: %w", err)
	}
	return nil
}

func (s *CartService) AddToCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	cart, err := s.getActiveCart(ctx, userID)
	// If err != nil then there is no cart, create one
	if err != nil {
		cart = models.Cart{
			ID:     uuid.New(),
			UserID: userID,
			Items:  []models.CartItem{},
		}

		err = s.db.CreateCart(ctx, database.CreateCartParams{
			ID:     cart.ID,
			UserID: cart.UserID,
		})
		if err != nil {
			return err
		}
	}

	// Add item to cart
	err = s.db.AddItemToCart(ctx, database.AddItemToCartParams{
		ID:        uuid.New(),
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  int32(quantity),
	})
	if err != nil {
		s.logger.Error("failed to add item to cart", zap.Error(err))
		return err
	}

	return nil
}

func (s *CartService) ReduceFromCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	// Get active cart
	cart, err := s.getActiveCart(ctx, userID)
	if err != nil {
		return err
	}

	// Reduce item from cart
	err = s.db.ReduceItemFromCart(ctx, database.ReduceItemFromCartParams{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  int32(quantity),
	})
	if err == nil {
		return nil
	}

	if apperrors.IsCheckViolation(err) {
		err = s.db.RemoveItemFromCart(ctx, database.RemoveItemFromCartParams{
			CartID:    cart.ID,
			ProductID: productID,
		})
		if err != nil {
			return fmt.Errorf("failed to remove item from cart: %w", err)
		}
		return nil
	}
	s.logger.Error("failed to reduce item from cart", zap.Error(err))
	return fmt.Errorf("failed to reduce item from cart: %w", err)
}

func (s *CartService) RemoveFromCart(ctx context.Context, userID, productID uuid.UUID) error {
	cart, err := s.getActiveCart(ctx, userID)
	if err != nil {
		return err
	}
	return s.db.RemoveItemFromCart(ctx, database.RemoveItemFromCartParams{
		CartID:    cart.ID,
		ProductID: productID,
	})
}

func (s *CartService) getActiveCart(ctx context.Context, userID uuid.UUID) (models.Cart, error) {
	// Get active cart
	cartRecord, err := s.db.GetActiveCart(ctx, userID)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			return models.Cart{}, fmt.Errorf("no active cart found for user: %w", apperrors.ErrNotFound)
		}
		return models.Cart{}, fmt.Errorf("could not fetch user cart: %w", err)
	}

	cart := models.Cart{
		ID:     cartRecord.ID,
		UserID: userID,
		Status: cartRecord.Status,
	}

	// Fetch cart items
	cartItems, err := s.db.GetCartItems(ctx, cart.ID)
	if err != nil {
		return models.Cart{}, err
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
