package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/synapsis-test/internal/user"
)

var (
	errUnknownConfig = errors.New("unknown config name")
)

// Handler contains user HTTP-handlers.
type Handler struct {
	handlers map[string]*handler
	user     user.Service
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
	// HandlerUsers denotes HTTP handler to interact
	// with users
	HandlerUsers = HandlerIdentity{
		Name: "users",
		URL:  "/v1/user",
	}

	// HandlerLogin denotes HTTP handler for user to login.
	HandlerLogin = HandlerIdentity{
		Name: "login",
		URL:  "/v1/user/login",
	}

	// HandlerPassword denotes HTTP handler for user to
	// interact with user password data.
	HandlerPassword = HandlerIdentity{
		Name: "password",
		URL:  "/v1/user/password/{id}",
	}

	// HandlerToken denotes HTTP handler for user to interact
	// with access token.
	HandlerToken = HandlerIdentity{
		Name: "token",
		URL:  "/v1/user/token/{id}",
	}
)

// New creates a new Handler.
func New(user user.Service, identities []HandlerIdentity) (*Handler, error) {
	h := &Handler{
		handlers: make(map[string]*handler),
		user:     user,
	}

	// apply options
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
	case HandlerUsers.Name:
		httpHandler = &usersHandler{
			user: h.user,
		}
	case HandlerLogin.Name:
		httpHandler = &loginHandler{
			user: h.user,
		}
	case HandlerPassword.Name:
		httpHandler = &passwordHandler{
			user: h.user,
		}
	case HandlerToken.Name:
		httpHandler = &tokenHandler{
			user: h.user,
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

type userHTTP struct {
	ID          *int64  `json:"id"`
	Email       *string `json:"email"`
	Name        *string `json:"name"`
	Password    *string `json:"password"`
	PhoneNumber *string `json:"phone_number"`
}
