package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/synapsis-test/global/helper"
	"github.com/synapsis-test/internal/order"
	"github.com/synapsis-test/internal/user"
)

type orderHandler struct {
	order  order.Service
	client user.Service
}

func (h *orderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[Order HTTP][orderHandler] Failed to parse order ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidOrderID.Error()})
		return
	}

	// handle based on HTTP request method
	switch r.Method {
	case http.MethodGet:
		h.handleGetOrderByID(w, r, orderID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *orderHandler) handleGetOrderByID(w http.ResponseWriter, r *http.Request, orderID int64) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 1000*time.Millisecond)
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
			log.Printf("[Order HTTP][handleGetOrderByID] Failed to get order by ID. orderID: %d, Err: %s\n", orderID, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan order.Order, 1)
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
		err = checkAccessToken(ctx, h.client, token, "handleGetOrderByID")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// TODO: add authorization flow with roles

		res, err := h.order.GetOrderByID(ctx, orderID)
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
				log.Printf("[Order HTTP][handleGetOrderByID] Internal error from GetOrderByID. orderID: %d. Err: %s\n", orderID, err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- res
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case res := <-resChan:
		// format order
		var o orderHTTP
		o, err = formatOrder(res)
		if err != nil {
			return
		}
		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: o,
		})
	}
}
