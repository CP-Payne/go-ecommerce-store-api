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

var paymentSource = &paypal.PaymentSource{
	Paypal: &paypal.PaymentSourcePaypal{
		ExperienceContext: paypal.PaymentSourcePaypalExperienceContext{
			BrandName:               "Niche Store",
			ShippingPreference:      "NO_SHIPPING",
			LandingPage:             "NO_PREFERENCE",
			UserAction:              "PAY_NOW",
			PaymentMethodPreference: "UNRESTRICTED",
			Locale:                  "en-US",
			ReturnURL:               "http://localhost:3000/payment/capture-order",
			CancelURL:               "http://localhost:3000/payment/cancel-order",
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

	logger := p.logger.With(
		zap.String("method", "CreateProcessorOrder"),
	)

	items, err := p.orderItemsToPaypalItems(order.OrderItems)
	if err != nil {
		logger.Error("failed to convert order items into paypal items", zap.Error(err))
		return nil, fmt.Errorf("failed to process order items into paypal items: %w", err)
	}

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

	processorOrder, err := p.client.CreateOrder(ctx, paypal.OrderIntentCapture, units, paymentSource, nil)
	if err != nil {
		logger.Error("failed to create paypal order", zap.Error(err), zap.String("orderID", order.ID.String()))
		return nil, fmt.Errorf("failed to create paypal order: %w", err)
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
	logger.Info("paypal order created", zap.String("orderID", orderResult.ID))
	return &orderResult, nil
}

func (p *PayPalProcessor) CaptureOrder(ctx context.Context, processorOrderID string) (*models.OrderResult, error) {

	logger := p.logger.With(
		zap.String("method", "CaptureOrder"),
		zap.String("paypalOrderID", processorOrderID),
	)

	orderResponse, err := p.client.CaptureOrder(ctx, processorOrderID, paypal.CaptureOrderRequest{
		PaymentSource: nil,
	})
	if err != nil {
		logger.Error("failed to capture paypal order", zap.Error(err))
		return &models.OrderResult{}, fmt.Errorf("failed to capture paypal order: %w", err)
	}

	orderResult := models.OrderResult{
		ID:           orderResponse.ID,
		ApproveLink:  "",
		Status:       orderResponse.Status,
		PaymentEmail: orderResponse.Payer.EmailAddress,
		PayerID:      orderResponse.Payer.PayerID,
	}

	logger.Info("paypal payment captured", zap.Any("orderResult", orderResult))
	return &orderResult, nil
}

func (p *PayPalProcessor) orderItemsToPaypalItems(orderItems []models.OrderItem) ([]paypal.Item, error) {
	if len(orderItems) <= 0 {
		return nil, fmt.Errorf("cannot conver orderItems to paypal items: %w", errors.New("no order items"))
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
