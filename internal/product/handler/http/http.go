package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/synapsis-test/internal/product"
	"github.com/synapsis-test/internal/user"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)

// Handler contains product HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	product  product.Service
	client   user.Service
}

// handler is the HTTP handler wrapper.
type handler struct {
	h        http.Handler
	identity HandlerIdentity
}

// HandlerIdentity denotes the identity of an HTTP hanlder.
type HandlerIdentity struct {
	Name string
	URL  string
}

// Followings are the known HTTP handler identities
var (
	// HandlerProducts denotes HTTP handler to interact
	// with products
	HandlerProducts = HandlerIdentity{
		Name: "products",
		URL:  "/v1/products",
	}

	// HandlerProduct denotes HTTP handler to interact
	// with a product.
	HandlerProduct = HandlerIdentity{
		Name: "product",
		URL:  "/v1/products/{id}",
	}

	// HandlerProductCart denotes HTTP handler to interact
	// with a product cart.
	HandlerProductCart = HandlerIdentity{
		Name: "product-cart",
		URL:  "/v1/products/{id}/carts",
	}

	// HandlerProductsCart denotes HTTP handler to interact
	// with products cart.
	HandlerProductsCart = HandlerIdentity{
		Name: "products-cart",
		URL:  "/v1/products/carts/{id}",
	}
)

// New creates a new Handler.
func New(product product.Service, client user.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		product:  product,
		client:   client,
	}

	// apply identity
	for _, identity := range identities {
		if h.handlers == nil {
			h.handlers = map[string]*handler{}
		}

		h.handlers[identity.Name] = &handler{
			identity: identity,
		}

		handler, err := h.createHTTPHandler(identity.Name)
		if err != nil {
			return nil, err
		}

		h.handlers[identity.Name].h = handler
	}

	return h, nil
}

// createHTTPHandler creates a new HTTP handler that
// implements http.Handler.
func (h *Handler) createHTTPHandler(configName string) (http.Handler, error) {
	var httpHandler http.Handler
	switch configName {
	case HandlerProducts.Name:
		httpHandler = &productsHandler{
			product: h.product,
			client:  h.client,
		}
	case HandlerProduct.Name:
		httpHandler = &productHandler{
			product: h.product,
			client:  h.client,
		}
	case HandlerProductsCart.Name:
		httpHandler = &productsCartHandler{
			product: h.product,
			client:  h.client,
		}
	case HandlerProductCart.Name:
		httpHandler = &productCartHandler{
			product: h.product,
			client:  h.client,
		}
	default:
		return httpHandler, errUnknownConfig
	}
	return httpHandler, nil
}

// Start starts all HTTP handlers.
func (h *Handler) Start(multiplexer *mux.Router) error {
	for _, handler := range h.handlers {
		multiplexer.Handle(handler.identity.URL, handler.h)
	}
	return nil
}

type productHTTP struct {
	ID           *int64  `json:"id"`
	Name         *string `json:"name"`
	Price        *int64  `json:"price"`
	Description  *string `json:"description"`
	CategoryID   *int64  `json:"category_id"`
	CategoryName *string `json:"category_name"`
}

type cartHTTP struct {
	UserID       *int64  `json:"user_id"`
	ProductID    *int64  `json:"product_id"`
	ProductName  *string `json:"product_name"`
	ProductPrice *int64  `json:"product_price"`
	Quantity     *int64  `json:"quantity"`
}
