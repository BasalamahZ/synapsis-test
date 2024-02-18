package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/synapsis-test/global/helper"
	"github.com/synapsis-test/internal/order"
	"github.com/synapsis-test/internal/user"
)

type ordersHandler struct {
	order  order.Service
	client user.Service
}

func (h *ordersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreateOrder(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *ordersHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("[Order HTTP][handleCreateOrder] Failed to create order, Err: %s\n", err.Error())
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
		request := orderHTTP{}
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// format HTTP request into service object
		reqOrder, err := parseOrderFromCreateRequest(request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		orderID, err := h.order.CreateOrder(ctx, reqOrder)
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
				log.Printf("[Order HTTP][handleCreateOrder] Internal error from CreateOrder. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- orderID
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case orderID := <-resChan:
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: orderID,
		})
	}
}

// parseOrderFromCreateRequest returns Order from the
// given HTTP request object.
func parseOrderFromCreateRequest(req orderHTTP) (order.Order, error) {
	result := order.Order{}

	if req.UserID != nil {
		result.UserID = *req.UserID
	}

	if req.ProductID != nil {
		result.ProductID = *req.ProductID
	}

	if req.Quantity != nil {
		result.Quantity = *req.Quantity
	}

	return result, nil
}
