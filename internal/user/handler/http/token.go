package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/synapsis-test/global/helper"
	"github.com/synapsis-test/internal/user"
	"github.com/gorilla/mux"
)

type tokenHandler struct {
	user user.Service
}

func (h *tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[User HTTP][tokenHandler] Failed to parse user ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidUserID.Error()})
		return
	}

	// handle based on HTTP request method
	switch r.Method {
	case http.MethodPost:
		h.handleRefreshToken(w, r, userID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *tokenHandler) handleRefreshToken(w http.ResponseWriter, r *http.Request, userID int64) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Millisecond)
	defer cancel()

	var (
		err        error           // stores error in this handler
		source     string          // stores request source
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		if err != nil {
			log.Printf("[User HTTP][handleRefreshToken] Failed to refresh token. Source: %s, Err: %s\n", source, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		// get token from header
		token, err := helper.GetBearerTokenFromHeader(r)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errInvalidToken
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.user, token, "handleRefreshToken", userID)
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// TODO: add authorization flow with roles

		// refresh token
		newToken, err := h.user.RefreshToken(ctx, token)
		if err != nil {
			// determine error and status code, by default its internal error
			parsedErr := errInternalServer
			statusCode = http.StatusInternalServerError
			if v, ok := mapHTTPError[err]; ok {
				parsedErr = v
				statusCode = http.StatusBadRequest
			}

			// log the actual error if its internal error
			if statusCode == http.StatusInternalServerError {
				log.Printf("[User HTTP][handleRefreshToken] Internal error from RefreshToken. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- newToken
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case token := <-resChan:
		res := helper.ResponseEnvelope{
			Data: token,
		}
		resBody, err = json.Marshal(res)
	}
}
