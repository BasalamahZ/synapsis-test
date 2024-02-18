package http

import (
	"errors"

	"github.com/synapsis-test/internal/product"
)

// Followings are the known errors from Product HTTP handlers.
var (
	// errDataNotFound is returned when the desired data is
	// not found.
	errDataNotFound = errors.New("DATA_NOT_FOUND")

	// errBadRequest is returned when the given request is
	// bad/invalid.
	errBadRequest = errors.New("BAD_REQUEST")

	// errInternalServer is returned when there is an
	// unexpected error encountered when processing a request.
	errInternalServer = errors.New("INTERNAL_SERVER_ERROR")

	// errMethodNotAllowed is returned when accessing not
	// allowed HTTP method.
	errMethodNotAllowed = errors.New("METHOD_NOT_ALLOWED")

	// errRequestTimeout is returned when processing time has
	// reached the timeout limit.
	errRequestTimeout = errors.New("REQUEST_TIMEOUT")

	// errInvalidProductID is returned when the given product ID is
	// invalid.
	errInvalidProductID = errors.New("INVALID_PRODUCT_ID")

	// errInvalidUserID is returned when the given user ID is
	// invalid.
	errInvalidUserID = errors.New("INVALID_USER_ID")

	// errInvalidQuantity is returned when the given quantity is
	// invalid.
	errInvalidQuantity = errors.New("INVALID_QUANTITY")

	// errInvalidToken is returned when the given token is
	// invalid.
	errInvalidToken = errors.New("INVALID_TOKEN")

	// errUnauthorizedAccess is returned when the request
	// is unaothorized.
	errUnauthorizedAccess = errors.New("UNAUTHORIZED_ACCESS")
)

var (
	// mapHTTPError maps service error into HTTP error that
	// categorize as bad request error.
	//
	// Internal server error-related should not be mapped here,
	// and the handler should just return `errInternal` as the
	// error instead
	mapHTTPError = map[error]error{
		product.ErrDataNotFound:     errDataNotFound,
		product.ErrInvalidProductID: errInvalidProductID,
		product.ErrInvalidUserID:    errInvalidUserID,
		product.ErrInvalidQuantity:    errInvalidQuantity,
	}
)
