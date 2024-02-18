package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/synapsis-test/internal/order"
)

func (s *service) payment(req order.Order) (*coreapi.ChargeResponse, error) {
	midtrans.ServerKey = "SB-Mid-server-6qHe2NyZzS7qYRI3Qskechx5"
	midtrans.Environment = midtrans.Sandbox

	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(1e10)
	orderID := fmt.Sprintf("%010d", randomNum)
	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeQris,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: req.TotalAmount,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    fmt.Sprintf("%v", req.ProductID),
				Name:  req.ProductName,
				Price: req.ProductPrice,
				Qty:   int32(req.Quantity),
			},
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: req.UserName,
			Email: req.UserEmail,
		},
		CustomExpiry: &coreapi.CustomExpiry{
			ExpiryDuration: 1,
			Unit:           "day",
		},
	}

	coreApiRes, err := coreapi.ChargeTransaction(chargeReq)
	if err != nil {
		return nil, err
	}

	return coreApiRes, nil
}

func (s *service) checkStatusPayment(orderID string) (*coreapi.TransactionStatusResponse, error) {
	res, err := coreapi.CheckTransaction(orderID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
