package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/synapsis-test/internal/category"
	categoryhttphandler "github.com/synapsis-test/internal/category/handler/http"
	categoryservice "github.com/synapsis-test/internal/category/service"
	categorypgstore "github.com/synapsis-test/internal/category/store/postgresql"
	"github.com/synapsis-test/internal/order"
	orderhttphandler "github.com/synapsis-test/internal/order/handler/http"
	orderservice "github.com/synapsis-test/internal/order/service"
	orderpgstore "github.com/synapsis-test/internal/order/store/postgresql"
	"github.com/synapsis-test/internal/product"
	producthttphandler "github.com/synapsis-test/internal/product/handler/http"
	productservice "github.com/synapsis-test/internal/product/service"
	productpgstore "github.com/synapsis-test/internal/product/store/postgresql"
	"github.com/synapsis-test/internal/user"
	userhttphandler "github.com/synapsis-test/internal/user/handler/http"
	userservice "github.com/synapsis-test/internal/user/service"
	userpgstore "github.com/synapsis-test/internal/user/store/postgresql"
)

// Following constants are the possible exit code returned
// when running a server.
const (
	CodeSuccess = iota
	CodeBadConfig
	CodeFailServeHTTP
)

// Run creates a server and starts the server.
//
// Run returns a status code suitable for os.Exit() argument.
func Run() int {
	s, err := new()
	if err != nil {
		return CodeBadConfig
	}

	return s.start()
}

// server is the long-runnning application.
type server struct {
	srv      *http.Server
	handlers []handler
}

// handler provides mechanism to start HTTP handler. All HTTP
// handlers must implements this interface.
type handler interface {
	Start(multiplexer *mux.Router) error
}

// new creates and returns a new server.
func new() (*server, error) {
	s := &server{
		srv: &http.Server{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	// connect to database
	db, err := sqlx.Connect("postgres", BaseConfig())
	if err != nil {
		log.Printf("[synapsistest-api-http] failed to connect database: %s\n", err.Error())
		return nil, fmt.Errorf("failed to connect database: %s", err.Error())
	}

	// connect to cache redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	// initialize user service
	var userSvc user.Service
	{
		pgStore, err := userpgstore.New(db)
		if err != nil {
			log.Printf("[user-api-http] failed to initialize user postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize user postgresql store: %s", err.Error())
		}

		svcOptions := []userservice.Option{}
		svcOptions = append(svcOptions, userservice.WithConfig(userservice.Config{
			PasswordSalt:   os.Getenv("PasswordSalt"),
			TokenSecretKey: os.Getenv("TokenSecretKey"),
		}))

		userSvc, err = userservice.New(pgStore, svcOptions...)
		if err != nil {
			log.Printf("[user-api-http] failed to initialize user service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize user service: %s", err.Error())
		}
	}

	// initialize product service
	var productSvc product.Service
	{
		pgStore, err := productpgstore.New(db)
		if err != nil {
			log.Printf("[product-api-http] failed to initialize product postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize product postgresql store: %s", err.Error())
		}

		productSvc, err = productservice.New(pgStore, rdb)
		if err != nil {
			log.Printf("[product-api-http] failed to initialize product service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize product service: %s", err.Error())
		}
	}

	// initialize category service
	var categorySvc category.Service
	{
		pgStore, err := categorypgstore.New(db)
		if err != nil {
			log.Printf("[category-api-http] failed to initialize category postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize category postgresql store: %s", err.Error())
		}

		categorySvc, err = categoryservice.New(pgStore)
		if err != nil {
			log.Printf("[category-api-http] failed to initialize category service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize category service: %s", err.Error())
		}
	}

	// initialize order service
	var orderSvc order.Service
	{
		pgStore, err := orderpgstore.New(db)
		if err != nil {
			log.Printf("[order-api-http] failed to initialize order postgresql store: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize order postgresql store: %s", err.Error())
		}

		orderSvc, err = orderservice.New(pgStore, productSvc)
		if err != nil {
			log.Printf("[order-api-http] failed to initialize order service: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize order service: %s", err.Error())
		}
	}

	// initialize user HTTP handler
	{
		identities := []userhttphandler.HandlerIdentity{
			userhttphandler.HandlerPassword,
			userhttphandler.HandlerToken,
			userhttphandler.HandlerLogin,
			userhttphandler.HandlerUsers,
		}

		userHTTP, err := userhttphandler.New(userSvc, identities)
		if err != nil {
			log.Printf("[user-api-http] failed to initialize user http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize user http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, userHTTP)
	}

	// initialize product HTTP handler
	{
		identities := []producthttphandler.HandlerIdentity{
			producthttphandler.HandlerProductsCart,
			producthttphandler.HandlerProduct,
			producthttphandler.HandlerProducts,
			producthttphandler.HandlerProductCart,
		}

		productHTTP, err := producthttphandler.New(productSvc, userSvc, identities)
		if err != nil {
			log.Printf("[product-api-http] failed to initialize product http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize product http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, productHTTP)
	}

	// initialize category HTTP handler
	{
		identities := []categoryhttphandler.HandlerIdentity{
			categoryhttphandler.HandlerCategory,
			categoryhttphandler.HandlerCategories,
		}

		categoryHTTP, err := categoryhttphandler.New(categorySvc, userSvc, identities)
		if err != nil {
			log.Printf("[category-api-http] failed to initialize category http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize category http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, categoryHTTP)
	}

	// initialize order HTTP handler
	{
		identities := []orderhttphandler.HandlerIdentity{
			orderhttphandler.HandlerOrder,
			orderhttphandler.HandlerOrders,
		}

		orderHTTP, err := orderhttphandler.New(orderSvc, userSvc, identities)
		if err != nil {
			log.Printf("[order-api-http] failed to initialize order http handlers: %s\n", err.Error())
			return nil, fmt.Errorf("failed to initialize order http handlers: %s", err.Error())
		}

		s.handlers = append(s.handlers, orderHTTP)
	}

	return s, nil
}

// start starts the given server.
func (s *server) start() int {
	log.Println("[synapsis-test-api-http] starting server...")

	// create multiplexer object
	rootMux := mux.NewRouter()
	appMux := rootMux.PathPrefix("/api").Subrouter()

	// starts handlers
	for _, h := range s.handlers {
		if err := h.Start(appMux); err != nil {
			log.Printf("[synapsis-test-api-http] failed to start handler: %s\n", err.Error())
			return CodeFailServeHTTP
		}
	}

	// endpoint checker
	appMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello world! @synapsis-test")
	})

	// use middlewares to app mux only
	appMux.Use(corsMiddleware)

	// listen and serve
	log.Printf("[synapsis-test-api-http] Server is running at %s:%s", os.Getenv("ADDRESS"), os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("ADDRESS"), os.Getenv("PORT")), rootMux))

	return CodeSuccess
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
