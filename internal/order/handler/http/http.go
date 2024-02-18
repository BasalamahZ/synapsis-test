package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/synapsis-test/internal/order"
	"github.com/synapsis-test/internal/user"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)

// Handler contains order HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	order    order.Service
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
	// HandlerOrders denotes HTTP handler to interact
	// with orders
	HandlerOrders = HandlerIdentity{
		Name: "orders",
		URL:  "/v1/orders",
	}

	// HandlerOrder denotes HTTP handler to interact
	// with a order.
	HandlerOrder = HandlerIdentity{
		Name: "order",
		URL:  "/v1/orders/{id}",
	}
)

// New creates a new Handler.
func New(order order.Service, client user.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		order:    order,
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
	case HandlerOrders.Name:
		httpHandler = &ordersHandler{
			order:  h.order,
			client: h.client,
		}
	case HandlerOrder.Name:
		httpHandler = &orderHandler{
			order:  h.order,
			client: h.client,
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

type orderHTTP struct {
	ID               *int64           `json:"id"`
	UserID           *int64           `json:"user_id"`
	UserName         *string          `json:"user_name"`
	UserEmail        *string          `json:"user_email"`
	ProductID        *int64           `json:"product_id"`
	ProductName      *string          `json:"product_name"`
	ProductPrice     *int64           `json:"product_price"`
	Quantity         *int64           `json:"quantity"`
	TotalAmount      *int64           `json:"total_amount"`
	Status           *string          `json:"status"`
	ResponseMidtrans *json.RawMessage `json:"response_midtrans"`
}
