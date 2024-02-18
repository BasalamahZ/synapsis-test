package http

import (
	"context"
	"log"

	"github.com/synapsis-test/internal/user"
)

// checkAccessToken checks the given access token whether it
// is valid or not.
func checkAccessToken(ctx context.Context, svc user.Service, token, name string, userID int64) error {
	tokenData, err := svc.ValidateToken(ctx, token)
	if err != nil {
		parsedErr := errUnauthorizedAccess
		if v, ok := mapHTTPError[err]; ok {
			parsedErr = v
		} else {
			log.Printf("[User HTTP][%s] Unauthorized error from ValidateToken. Err: %s\n", name, err.Error())
		}

		return parsedErr
	}

	if userID != tokenData.UserID {
		return errInvalidUserID
	}

	return nil
}
