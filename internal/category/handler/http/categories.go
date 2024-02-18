package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/synapsis-test/global/helper"
	"github.com/synapsis-test/internal/category"
	"github.com/synapsis-test/internal/user"
)

type categoriesHandler struct {
	category category.Service
	client   user.Service
}

func (h *categoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetCategories(w, r)
	default:
		helper.WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{errMethodNotAllowed.Error()})
	}
}

func (h *categoriesHandler) handleGetCategories(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("[Category HTTP][handleGetCategories] Failed to get all categories. Err: %s\n", err.Error())
			helper.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		helper.WriteResponse(w, resBody, statusCode, helper.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan []category.Category, 1)
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
		err = checkAccessToken(ctx, h.client, token, "handleGetCategories")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		res, err := h.category.GetCategories(ctx)
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
				log.Printf("[Category HTTP][handleGetCategories] Internal error from GetCategories. Err: %s\n", err.Error())
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
		// format each categories
		categories := make([]categoryHTTP, 0)
		for _, r := range res {
			var c categoryHTTP
			c, err = formatCategory(r)
			if err != nil {
				return
			}
			categories = append(categories, c)
		}

		// construct response data
		resBody, err = json.Marshal(helper.ResponseEnvelope{
			Data: categories,
		})
	}
}
