package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/synapsis-test/internal/category"
	"github.com/synapsis-test/internal/user"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)


// Handler contains category HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	category  category.Service
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
	// HandlerCategories denotes HTTP handler to interact
	// with categories
	HandlerCategories = HandlerIdentity{
		Name: "categories",
		URL:  "/v1/categories",
	}

	// HandlerCategory denotes HTTP handler to interact
	// with a category.
	HandlerCategory = HandlerIdentity{
		Name: "category",
		URL:  "/v1/categories/{id}",
	}
)

// New creates a new Handler.
func New(category category.Service, client user.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		category:  category,
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
	case HandlerCategories.Name:
		httpHandler = &categoriesHandler{
			category: h.category,
			client:  h.client,
		}
	case HandlerCategory.Name:
		httpHandler = &categoryHandler{
			category: h.category,
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

type categoryHTTP struct {
	ID           *int64  `json:"id"`
	Name         *string `json:"name"`
	Description  *string `json:"description"`
}
