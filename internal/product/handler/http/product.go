package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/synapsis-test/global/helper"
	"github.com/synapsis-test/internal/product"
	"github.com/gorilla/mux"
	"github.com/synapsis-test/internal/user"
)

type productHandler struct {
	product product.Service
	client  user.Service
}

func (h *productHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[Product HTTP][productHandler] Failed to parse product ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidProductID.Error()})
		return
	}

	// handle based on HTTP request method
	switch r.Method {
	case http.MethodGet:
		h.handleGetProductByID(w, r, productID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *productHandler) handleGetProductByID(w http.ResponseWriter, r *http.Request, productID int64) {
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
			log.Printf("[Product HTTP][handleGetProductByID] Failed to get product by ID. productID: %d, Err: %s\n", productID, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan product.Product, 1)
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
		err = checkAccessToken(ctx, h.client, token, "handleGetProductByID")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// TODO: add authorization flow with roles

		res, err := h.product.GetProductByID(ctx, productID)
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
				log.Printf("[Product HTTP][handleGetProductByID] Internal error from GetProductByID. productID: %d. Err: %s\n", productID, err.Error())
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
		// format product
		var p productHTTP
		p, err = formatProduct(res)
		if err != nil {
			return
		}
		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: p,
		})
	}
}
