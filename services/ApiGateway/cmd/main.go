package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kareemhamed001/e-commerce/pkg/logger"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/config"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/internal/clients"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/internal/handlers"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/internal/router"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Errorf("failed to load configuration: %v", err)
		return
	}

	// Initialize logger
	logger.InitGlobal(cfg.AppEnv)
	logger.Info("Starting API Gateway...")

	logger.Infof("Configuration loaded successfully")

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize gRPC clients
	serviceClients, err := clients.NewServiceClients(
		cfg.UserServiceURL,
		cfg.ProductServiceURL,
		cfg.CartServiceURL,
		cfg.OrderServiceURL,
	)
	if err != nil {
		logger.Errorf("Failed to initialize service clients: %v", err)
		return
	}
	defer serviceClients.Close()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(serviceClients.UserClient)
	productHandler := handlers.NewProductHandler(serviceClients.ProductClient)
	cartHandler := handlers.NewCartHandler(serviceClients.CartClient)
	orderHandler := handlers.NewOrderHandler(serviceClients.OrderClient)

	routerEngine := gin.Default()

	// Initialize router
	apiRouter := router.NewRouter(routerEngine, cfg, userHandler, productHandler, cartHandler, orderHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      apiRouter.Handler(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("API Gateway listening on port %s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Failed to start server: %v", err)
		}
		logger.Info("Server stopped")

	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down API Gateway...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("API Gateway stopped")
}
