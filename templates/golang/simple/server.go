package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type APIServer struct {
	config     *Config
	logger     *zap.Logger
	router     http.Handler
	server     *http.Server
	httpClient *httpClient
}

func init() {
	time.LoadLocation("UTC")             // ensure all time is in UTC
	runtime.GOMAXPROCS(runtime.NumCPU()) // set the core
}

// Creates zap logger according to the config
func NewLogger(config *Config) (*zap.Logger, error) {
	c := zap.NewProductionConfig()
	c.DisableCaller = true

	if config.Verbose {
		c.DisableCaller = false
		c.Development = true
		c.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	return c.Build()
}

// NewAPIServer create's a new APIServer from configuration
func NewAPIServer(config *Config) (*APIServer, error) {
	// create the service logger
	logger, err := NewLogger(config)
	if err != nil {
		return nil, err
	}
	logger.Info("Starting the service",
		zap.String("prog", prog),
		zap.String("version", version))

	httpClient, httpErr := NewHttpClient(logger, config)
	if httpErr != nil {
		return nil, httpErr
	}

	svr := &APIServer{
		httpClient: httpClient,
		config:     config,
		logger:     logger,
	}
	svr.setupRoutes()

	return svr, nil
}

func (a *APIServer) setupRoutes() error {
	router := chi.NewRouter()
	router.MethodNotAllowed(MethodNotAllowed)
	router.NotFound(NotFound)
	router.Use(middleware.RequestID)
	router.Use(a.LoggingMiddleware)
	router.Use(a.IdentityMiddleware)
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Mount(fmt.Sprintf("/%s", prog), a.Routes())

	a.router = router

	return nil
}

func (a *APIServer) Run() error {
	a.logger.Info("groupsets API HTTP service starting", zap.Object("config", a.config))
	server := &http.Server{
		Addr:    a.config.Addr,
		Handler: a.router,
	}
	a.server = server
	err := a.Init()
	if err != nil {
		return err
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			a.logger.Fatal("failed to start the API HTTP server", zap.Error(err))
		}
	}()
	return nil
}

// Initialize api server, called before running the server
func (a *APIServer) Init() error {
	return nil
}
