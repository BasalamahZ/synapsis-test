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

type loginHandler struct {
	user user.Service
}

// loginRequestData is the data from user to perform login.
type loginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginResponseData is the data to user after perform login.
type loginResponseData struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

func (h *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle based on HTTP request method
	switch r.Method {
	case http.MethodPost:
		h.handleLogin(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *loginHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Millisecond)
	defer cancel()

	var (
		err        error           // stores error in this handler
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		// error
		if err != nil {
			log.Printf("[User HTTP][handleLogin] Failed to login. Err: %s\n", err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan loginResponseData, 1)
	errChan := make(chan error, 1)

	go func() {
		// read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// unmarshall body
		var data loginRequestData
		err = json.Unmarshal(body, &data)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// login
		token, tokenData, err := h.user.LoginBasic(ctx, data.Email, data.Password)
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
				log.Printf("[User HTTP][handleLogin] Internal error from LoginBasic. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- loginResponseData{
			UserID: tokenData.UserID,
			Email:  tokenData.Email,
			Token:  token,
		}
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case resData := <-resChan:
		res := helper.ResponseEnvelope{
			Data: resData,
		}
		resBody, err = json.Marshal(res)
	}
}
