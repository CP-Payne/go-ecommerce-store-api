package service

import (
	"context"
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

// TODO: Hide the below functions???
func (p *PayPalProcessor) CreateProductOrder(ctx context.Context, product *models.Product, quantity int, shippingPrice float32) (*models.OrderResult, error) {
	total := product.Price*float32(quantity) + shippingPrice
	totalStr := fmt.Sprintf("%.2f", total)

	itemPrice := product.Price * float32(quantity)
	itemPriceStr := fmt.Sprintf("%.2f", itemPrice)
	shippingPriceStr := fmt.Sprintf("%.2f", shippingPrice)
	// TODO: Create separate function for price calculation and item creation
	// TODO: Add item information to the below struct
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
