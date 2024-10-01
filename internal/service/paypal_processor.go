package service

import (
	"context"
	"errors"
	"fmt"

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

// TODO: What happens if cart is empty and payment is made?

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
	// client.SetLog(os.Stdout)

	return &PayPalProcessor{
		logger:          logger,
		client:          client,
		processorConfig: pconf,
	}, nil
}

func (p *PayPalProcessor) CreateProcessorOrder(ctx context.Context, order *models.Order) (*models.OrderResult, error) {
	items, err := p.orderItemsToPaypalItems(order.OrderItems)
	if err != nil {
		p.logger.Error("failed to create order", zap.Error(err))
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// p.logger.Debug("items", zap.Any("items", items))
	units := []paypal.PurchaseUnitRequest{
		{
			Amount: &paypal.PurchaseUnitAmount{
				Currency: "USD",
				Value:    floatToString(order.OrderTotal),
				Breakdown: &paypal.PurchaseUnitAmountBreakdown{
					ItemTotal: &paypal.Money{
						Currency: "USD",
						Value:    floatToString(order.ProductTotal),
					},
					Shipping: &paypal.Money{
						Currency: "USD",
						Value:    floatToString(order.ShippingPrice),
					},
				},
			},
			Items: items,
		},
	}

	// p.logger.Debug("full units", zap.Any("units", units))

	processorOrder, err := p.client.CreateOrder(ctx, paypal.OrderIntentCapture, units, paymentSource, nil)
	if err != nil {
		p.logger.Error("failed to create order", zap.Error(err), zap.String("orderID", order.ID.String()))
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	orderLinks := processorOrder.Links
	var approveLink string

	for _, l := range orderLinks {
		if l.Rel == "payer-action" {
			approveLink = l.Href
		}
	}

	orderResult := models.OrderResult{
		ID:          processorOrder.ID,
		ApproveLink: approveLink,
		Status:      "",
	}
	return &orderResult, nil
}

// TODO: Fix dry code
// TODO: Create separate function for price calculation and item creation

func (p *PayPalProcessor) CaptureOrder(ctx context.Context, processorOrderID string) (*models.OrderResult, error) {
	// TODO: Return orderResponse and save it to db
	orderResponse, err := p.client.CaptureOrder(ctx, processorOrderID, paypal.CaptureOrderRequest{
		PaymentSource: nil,
	})
	if err != nil {
		p.logger.Error("failed to capture order", zap.Error(err))
		return &models.OrderResult{}, fmt.Errorf("failed to capture order: %w", err)
	}

	orderResult := models.OrderResult{
		ID:           orderResponse.ID,
		ApproveLink:  "",
		Status:       orderResponse.Status,
		PaymentEmail: orderResponse.Payer.EmailAddress,
		PayerID:      orderResponse.Payer.PayerID,
	}

	return &orderResult, nil
}

// TODO: Add the new error to apperrors

func (p *PayPalProcessor) orderItemsToPaypalItems(orderItems []models.OrderItem) ([]paypal.Item, error) {
	if len(orderItems) <= 0 {
		p.logger.Error("failed to convert orderItems to paypal items. no orderItems provided")
		return nil, fmt.Errorf("failed to convert orderItems to paypal items: %w", errors.New("no order items"))
	}

	paypalItems := make([]paypal.Item, 0, len(orderItems))

	for _, oi := range orderItems {
		paypalItem := paypal.Item{
			Name: oi.Name,
			UnitAmount: &paypal.Money{
				Currency: "USD",
				Value:    floatToString(oi.Price),
			},
			Quantity: intToString(oi.Quantity),
		}
		paypalItems = append(paypalItems, paypalItem)
	}

	return paypalItems, nil
}

// func (p *PayPalProcessor) GetOrderDetails(ctx context.Context, orderID string) error {
// 	order, err := p.client.GetOrder(ctx, orderID)
// 	if err != nil {
// 		p.logger.Error("failed to retrieve order details", zap.Error(err))
// 		return fmt.Errorf("failed to retrieve order details: %w", err)
// 	}
//
// 	p.logger.Debug("GetOrderDetails: ", zap.Any("ORDER", order))
//
// 	return nil
// }
