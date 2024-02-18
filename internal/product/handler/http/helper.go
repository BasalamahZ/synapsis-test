package http

import (
	"github.com/synapsis-test/internal/product"
)

// formatProduct formats the given product
// into the respective HTTP-format object.
func formatProduct(p product.Product) (productHTTP, error) {
	return productHTTP{
		ID:           &p.ID,
		Name:         &p.Name,
		Price:        &p.Price,
		Description:  &p.Description,
		CategoryID:   &p.CategoryID,
		CategoryName: &p.CategoryName,
	}, nil
}

// formatProductCarrt formats the given product cart
// into the respective HTTP-format object.
func formatProductCart(pc product.ProductCart) (cartHTTP, error) {
	return cartHTTP{
		UserID:       &pc.UserID,
		ProductID:    &pc.ProductID,
		ProductName:  &pc.ProductName,
		ProductPrice: &pc.ProductPrice,
		Quantity:     &pc.Quantity,
	}, nil
}
