package http

import (
	"encoding/json"

	"github.com/synapsis-test/internal/order"
)

// formatOrder formats the given order
// into the respective HTTP-format object.
func formatOrder(o order.Order) (orderHTTP, error) {
	statusStr := o.Status.String()
	// responseMidtransString := o.ResponseMidtrans
	// responseMidtransBytes := []byte(responseMidtransString)
	responseMidtrans := json.RawMessage([]byte(o.ResponseMidtrans))
	return orderHTTP{
		ID:               &o.ID,
		UserID:           &o.UserID,
		UserName:         &o.UserName,
		UserEmail:        &o.UserEmail,
		ProductID:        &o.ProductID,
		ProductName:      &o.ProductName,
		ProductPrice:     &o.ProductPrice,
		Quantity:         &o.Quantity,
		TotalAmount:      &o.TotalAmount,
		Status:           &statusStr,
		ResponseMidtrans: &responseMidtrans,
	}, nil
}
