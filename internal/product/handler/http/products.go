package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/synapsis-test/global/helper"
	"github.com/synapsis-test/internal/product"
	"github.com/synapsis-test/internal/user"
)

type productsHandler struct {
	product product.Service
	client  user.Service
}

func (h *productsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetProducts(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *productsHandler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	ctx, cancel := context.WithTimeout(r.Context(), 2000*time.Millisecond)
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
			log.Printf("[Product HTTP][handleGetProducts] Failed to get all Products. Err: %s\n", err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan []product.Product, 1)
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
		err = checkAccessToken(ctx, h.client, token, "handleGetProducts")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		// parsed filter
		filter, err := parseGetProductsFilter(r.URL.Query())
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- err
			return
		}

		res, err := h.product.GetProducts(ctx, filter)
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
				log.Printf("[Product HTTP][handleGetProducts] Internal error from GetProducts. Err: %s\n", err.Error())
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
		// format each products
		products := make([]productHTTP, 0)
		for _, r := range res {
			var p productHTTP
			p, err = formatProduct(r)
			if err != nil {
				return
			}
			products = append(products, p)
		}

		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: products,
		})
	}
}

func parseGetProductsFilter(request url.Values) (product.GetProductsFilter, error) {
	result := product.GetProductsFilter{}

	var categoryID int64
	if categoryIDStr := request.Get("category_id"); categoryIDStr != "" {
		intCategoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			return result, err
		}
		categoryID = intCategoryID
	}


	return product.GetProductsFilter{
		CategoryID:  categoryID,
	}, nil
}
