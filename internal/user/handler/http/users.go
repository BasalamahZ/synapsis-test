package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/synapsis-test/global/helper"
	"github.com/synapsis-test/internal/user"
)

type usersHandler struct {
	user user.Service
}

func (h *usersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreateUser(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *usersHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("[User HTTP][handleCreateUser] Failed to create user. Source: %s, Err: %s\n", source, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	go func() {
		// read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// unmarshall body
		request := userHTTP{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// format HTTP request into service object
		reqUser, err := parseUserFromCreateRequest(request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		userID, err := h.user.CreateUser(ctx, reqUser)
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
				log.Printf("[User HTTP][handleCreateUser] Internal error from CreateUser. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- userID
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case userID := <-resChan:
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Status: "Success",
			Data:   userID,
		})
	}
}

// parseUserFromCreateRequest returns user
// from the given HTTP request object.
func parseUserFromCreateRequest(uh userHTTP) (user.User, error) {
	result := user.User{}

	if uh.Email != nil {
		result.Email = *uh.Email
	}

	if uh.Name != nil {
		result.Name = *uh.Name
	}

	if uh.Password != nil {
		result.Password = *uh.Password
	}

	if uh.PhoneNumber != nil {
		result.PhoneNumber = *uh.PhoneNumber
	}

	return result, nil
}
