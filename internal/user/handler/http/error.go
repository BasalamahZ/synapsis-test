package http

import (
	"errors"

	"github.com/synapsis-test/internal/user"
)

// Followings are the known errors from User HTTP handlers.
var (
	// errUserAlreadyExist is returned when the given
	// user already exist based on the predefined
	// unique constraints.
	errUserAlreadyExist = errors.New("USER_ALREADY_EXIST")

	// errDataNotFound is returned when the desired data is
	// not found.
	errDataNotFound = errors.New("DATA_NOT_FOUND")

	// errBadRequest is returned when the given request is
	// bad/invalid.
	errBadRequest = errors.New("BAD_REQUEST")

	// errExpiredToken is returned when the given token is
	// expired.
	errExpiredToken = errors.New("EXPIRED_TOKEN")

	// errInternalServer is returned when there is an
	// unexpected error encountered when processing a request.
	errInternalServer = errors.New("INTERNAL_SERVER_ERROR")

	// errInvalidPassword is returned when the given password
	// is invalid.
	errInvalidPassword = errors.New("INVALID_PASSWORD")

	// errInvalidToken is returned when the given token is
	// invalid.
	errInvalidToken = errors.New("INVALID_TOKEN")

	// errInvalidEmail is returned when the given email
	// is invalid.
	errInvalidEmail = errors.New("INVALID_EMAIL")

	// errMethodNotAllowed is returned when accessing not
	// allowed HTTP method.
	errMethodNotAllowed = errors.New("METHOD_NOT_ALLOWED")

	// errRequestTimeout is returned when processing time has
	// reached the timeout limit.
	errRequestTimeout = errors.New("REQUEST_TIMEOUT")

	// errInvalidUserID is returned when the given user ID is
	// invalid.
	errInvalidUserID = errors.New("INVALID_USER_ID")

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
		user.ErrUserAlreadyExist: errUserAlreadyExist,
		user.ErrDataNotFound:     errDataNotFound,
		user.ErrInvalidEmail:     errInvalidEmail,
		user.ErrExpiredToken:     errExpiredToken,
		user.ErrInvalidPassword:  errInvalidPassword,
		user.ErrInvalidToken:     errInvalidToken,
	}
)
