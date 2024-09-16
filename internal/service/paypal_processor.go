package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/models"
	"github.com/plutov/paypal/v4"
	"go.uber.org/zap"
)

type PayPalProcessor struct {
	logger          *zap.Logger
	client          *paypal.Client
	processorConfig *config.ProcessorConfig
}

// TODO: Maybe add below in PayPalProcessor struct
var paymentSource = &paypal.PaymentSource{
	Paypal: &paypal.PaymentSourcePaypal{
		ExperienceContext: paypal.PaymentSourcePaypalExperienceContext{
			BrandName:               "Niche Store",
			ShippingPreference:      "NO_SHIPPING",
			LandingPage:             "NO_PREFERENCE",
			UserAction:              "PAY_NOW",
			PaymentMethodPreference: "UNRESTRICTED",
			Locale:                  "en-US",
			ReturnURL:               "http://localhost:3000/products/complete-order",
			CancelURL:               "http://localhost:3000/products/cancel-order",
		},
	},
}

func NewPayPalProcessor(pconf *config.ProcessorConfig) (*PayPalProcessor, error) {
	logger := config.GetLogger()
	client, err := paypal.NewClient(pconf.ClientID, pconf.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		logger.Error("failed to create paypal client", zap.Error(err))
		return nil, fmt.Errorf("failed to create paypal client: %w", err)
	}
	client.SetLog(os.Stdout)

	return &PayPalProcessor{
		logger:          logger,
		client:          client,
		processorConfig: pconf,
	}, nil
}

// TODO: On successfull purchase, clear users cart

// TODO: Fix dry code

func (p *PayPalProcessor) CreateCartOrder(ctx context.Context, cart *models.Cart, shippingPrice float32) (*models.OrderResult, error) {
	cartTotal, err := p.getCartTotal(cart)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart order")
	}
	cartTotalStr := fmt.Sprintf("%.2f", cartTotal)

	total := cartTotal + shippingPrice
	totalStr := fmt.Sprintf("%.2f", total)
	shippingPriceStr := fmt.Sprintf("%.2f", shippingPrice)

	items, err := p.cartToPaypalItems(cart)
	if err != nil {
		p.logger.Error("failed to create order", zap.Error(err))
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	p.logger.Debug("items", zap.Any("items", items))
	units := []paypal.PurchaseUnitRequest{
		{
			Amount: &paypal.PurchaseUnitAmount{
				Currency: "USD",
				Value:    totalStr,
				Breakdown: &paypal.PurchaseUnitAmountBreakdown{
					ItemTotal: &paypal.Money{
						Currency: "USD",
						Value:    cartTotalStr,
					},
					Shipping: &paypal.Money{
						Currency: "USD",
						Value:    shippingPriceStr,
					},
				},
			},
			Items: items,
		},
	}

	p.logger.Debug("full units", zap.Any("units", units))

	order, err := p.client.CreateOrder(ctx, paypal.OrderIntentCapture, units, paymentSource, nil)
	if err != nil {
		p.logger.Error("failed to create cart order", zap.Error(err), zap.String("cartID", cart.ID.String()))
		return nil, fmt.Errorf("failed to create cart order: %w", err)
	}
	orderLinks := order.Links
	var approveLink string

	for _, l := range orderLinks {
		if l.Rel == "payer-action" {
			approveLink = l.Href
		}
	}

	orderResult := models.OrderResult{
		ID:          order.ID,
		ApproveLink: approveLink,
		Status:      "",
	}
	return &orderResult, nil
}

// TODO: Hide the below functions???
func (p *PayPalProcessor) CreateProductOrder(ctx context.Context, product *models.Product, quantity int, shippingPrice float32) (*models.OrderResult, error) {
	total := product.Price*float32(quantity) + shippingPrice
	totalStr := fmt.Sprintf("%.2f", total)

	itemPrice := product.Price * float32(quantity)
	itemPriceStr := fmt.Sprintf("%.2f", itemPrice)
	shippingPriceStr := fmt.Sprintf("%.2f", shippingPrice)
	// TODO: Create separate function for price calculation and item creation
	// TODO: Add item information to the below struct
	items, err := p.productToPaypalItem(product, quantity)
	if err != nil {
		p.logger.Error("failed to create order", zap.Error(err))
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	units := []paypal.PurchaseUnitRequest{
		{
			Amount: &paypal.PurchaseUnitAmount{
				Currency: "USD",
				Value:    totalStr,
				Breakdown: &paypal.PurchaseUnitAmountBreakdown{
					ItemTotal: &paypal.Money{
						Currency: "USD",
						Value:    itemPriceStr,
					},
					Shipping: &paypal.Money{
						Currency: "USD",
						Value:    shippingPriceStr,
					},
				},
			},
			Items: items,
		},
	}

	order, err := p.client.CreateOrder(ctx, paypal.OrderIntentCapture, units, paymentSource, nil)
	if err != nil {
		p.logger.Error("failed to create product order", zap.Error(err), zap.String("productID", product.ID.String()))
		return nil, fmt.Errorf("failed to create product order: %w", err)
	}
	orderLinks := order.Links
	var approveLink string

	for _, l := range orderLinks {
		if l.Rel == "payer-action" {
			approveLink = l.Href
		}
	}

	orderResult := models.OrderResult{
		ID:          order.ID,
		ApproveLink: approveLink,
		Status:      "",
	}
	return &orderResult, nil
}

func (p *PayPalProcessor) CaptureOrder(ctx context.Context, orderID string) error {
	// TODO: Return orderResponse and save it to db
	_, err := p.client.CaptureOrder(ctx, orderID, paypal.CaptureOrderRequest{
		PaymentSource: nil,
	})
	if err != nil {
		p.logger.Error("failed to capture order", zap.Error(err))
		return fmt.Errorf("failed to capture order: %w", err)
	}

	return nil
}

// TODO: On successfull purchase, reduce stock

func (p *PayPalProcessor) cartToPaypalItems(cart *models.Cart) ([]paypal.Item, error) {
	if cart == nil {
		// TODO: Add the new error to apperrors
		p.logger.Error("failed to convert cart to paypal item. cart cannot be nil")
		return nil, fmt.Errorf("failed to convert cart to paypal item -> cart is nil: %w", errors.New("nil cart"))
	}

	if len(cart.Items) <= 0 {
		p.logger.Error("failed to convert cart to paypal item. cart has no items")
		return nil, fmt.Errorf("failed to convert cart to paypal item -> no items in cart: %w", errors.New("empty cart"))
	}

	paypalItems := make([]paypal.Item, 0, len(cart.Items))

	for _, ci := range cart.Items {
		itemPriceStr := fmt.Sprintf("%.2f", ci.Price)
		quantityStr := fmt.Sprintf("%d", ci.Quantity)
		paypalItem := paypal.Item{
			Name: ci.Name,
			UnitAmount: &paypal.Money{
				Currency: "USD",
				Value:    itemPriceStr,
			},
			Quantity: quantityStr,
		}
		paypalItems = append(paypalItems, paypalItem)
	}

	return paypalItems, nil
}

func (p *PayPalProcessor) productToPaypalItem(product *models.Product, quantity int) ([]paypal.Item, error) {
	if product == nil {
		// TODO: Add the new error to apperrors
		p.logger.Error("failed to convert product to paypal item. product cannot be nil")
		return nil, fmt.Errorf("failed to convert product to paypal item -> product is nil: %w", errors.New("nil product"))
	}

	paypalItems := make([]paypal.Item, 0, 1)

	itemPriceStr := fmt.Sprintf("%.2f", product.Price)
	quantityStr := fmt.Sprintf("%d", quantity)
	paypalItem := paypal.Item{
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.Sku,
		UnitAmount: &paypal.Money{
			Currency: "USD",
			Value:    itemPriceStr,
		},
		Quantity: quantityStr,
	}
	paypalItems = append(paypalItems, paypalItem)

	return paypalItems, nil
}

func (p *PayPalProcessor) getCartTotal(cart *models.Cart) (float32, error) {
	if cart == nil {
		p.logger.Error("failed to calculate cart total. cart cannot be nil")
		return 0, fmt.Errorf("failed to calculate cart total -> cart is nil: %w", errors.New("nil cart"))
	}
	var cartTotal float32 = 0
	for _, ci := range cart.Items {
		cartTotal += ci.Price * float32(ci.Quantity)
	}

	return cartTotal, nil
}
