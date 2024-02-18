package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/synapsis-test/global/helper"
	"github.com/synapsis-test/internal/product"
	"github.com/synapsis-test/internal/user"
)

type productCartHandler struct {
	product product.Service
	client  user.Service
}

func (h *productCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[Product Cart HTTP][productCartHandler] Failed to parse product ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidProductID.Error()})
		return
	}

	// handle based on HTTP request method
	switch r.Method {
	case http.MethodPost:
		h.handleAddProductCart(w, r, productID)
	case http.MethodDelete:
		h.handleDeleteProductCartByID(w, r, productID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *productCartHandler) handleAddProductCart(w http.ResponseWriter, r *http.Request, productID int64) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 3000*time.Millisecond)
	defer cancel()

	var (
		err        error           // stores error in this handler
		request    cartHTTP        // stores request
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		// error
		if err != nil {
			log.Printf("[Product Cart HTTP][handleAddProductCart] Failed to add product cart, Err: %s\n", err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan struct{}, 1)
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
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// format HTTP request into service object
		reqProductCart, err := parseProductCartFromCreateRequest(request, productID)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		err = h.product.AddProductCart(ctx, reqProductCart)
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
				log.Printf("[Product Cart HTTP][handleAddProductCart] Internal error from AddProductCart. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- struct{}{}
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case <-resChan:
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: map[string]interface{}{
				"user_id":    request.UserID,
				"product_id": productID,
			},
		})
	}
}

func (h *productCartHandler) handleDeleteProductCartByID(w http.ResponseWriter, r *http.Request, productID int64) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 2000*time.Millisecond)
	defer cancel()

	var (
		err        error           // stores error in this handler
		request    cartHTTP        // stores request
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		// error
		if err != nil {
			log.Printf("[Product Cart HTTP][handleDeleteProductCartByID] Failed to delete product cart by ID. productID: %d, Err: %s\n", productID, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan struct{}, 1)
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
		err = json.Unmarshal(body, &request)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// format HTTP request into service object
		reqProductCart, err := parseProductCartFromCreateRequest(request, productID)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		err = h.product.DeleteProductCartByID(ctx, reqProductCart)
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
				log.Printf("[Product Cart HTTP][handleDeleteProductCartByID] Internal error from DeleteProductCartByID. productID: %d. Err: %s\n", productID, err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- struct{}{}
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case <-resChan:
		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: productID,
		})
	}
}

// parseProductCartFromCreateRequest returns product cart from the
// given HTTP request object.
func parseProductCartFromCreateRequest(req cartHTTP, productID int64) (product.ProductCart, error) {
	result := product.ProductCart{
		ProductID: productID,
	}

	if req.UserID != nil {
		result.UserID = *req.UserID
	}

	if req.Quantity != nil {
		result.Quantity = *req.Quantity
	}

	return result, nil
}
