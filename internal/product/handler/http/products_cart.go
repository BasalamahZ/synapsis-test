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
	"github.com/synapsis-test/internal/product"
	"github.com/synapsis-test/internal/user"
)

type productsCartHandler struct {
	product product.Service
	client  user.Service
}

func (h *productsCartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Printf("[Product Cart HTTP][productsCartHandler] Failed to parse user ID. ID: %s. Err: %s\n", vars["id"], err.Error())
		helper.WriteErrorResponse(w, http.StatusBadRequest, []string{errInvalidUserID.Error()})
		return
	}

	// handle based on HTTP request method
	switch r.Method {
	case http.MethodGet:
		h.handleGetCartsByUserID(w, r, userID)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *productsCartHandler) handleGetCartsByUserID(w http.ResponseWriter, r *http.Request, userID int64) {
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
			log.Printf("[Product Cart HTTP][handleGetCartsByUserID] Failed to get product by ID. userID: %d, Err: %s\n", userID, err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan []product.ProductCart, 1)
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
		err = checkAccessToken(ctx, h.client, token, "handleGetCartsByUserID")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// TODO: add authorization flow with roles

		res, err := h.product.GetCartsByUserID(ctx, userID)
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
				log.Printf("[Product Cart HTTP][handleGetCartsByUserID] Internal error from GetCartsByUserID. userID: %d. Err: %s\n", userID, err.Error())
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
		// format each carts
		carts := make([]cartHTTP, 0)
		for _, r := range res {
			var c cartHTTP
			c, err = formatProductCart(r)
			if err != nil {
				return
			}
			carts = append(carts, c)
		}
		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: carts,
		})
	}
}
