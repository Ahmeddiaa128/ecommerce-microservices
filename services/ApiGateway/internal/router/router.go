package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	customJWT "github.com/kareemhamed001/e-commerce/pkg/jwt"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/config"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/internal/handlers"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/internal/middleware"
)

// Router manages all HTTP routes and middlewares
type Router struct {
	engine         *gin.Engine
	cfg            *config.Config
	jwtManager     *customJWT.JWTManager
	userHandler    *handlers.UserHandler
	productHandler *handlers.ProductHandler
	cartHandler    *handlers.CartHandler
	orderHandler   *handlers.OrderHandler
}

// NewRouter creates a new router with all routes configured
func NewRouter(
	router *gin.Engine,
	cfg *config.Config,
	userHandler *handlers.UserHandler,
	productHandler *handlers.ProductHandler,
	cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler,
) *Router {
	r := &Router{
		engine:         router,
		cfg:            cfg,
		jwtManager:     customJWT.NewJWTManager(cfg.JWTSecret, 24*time.Hour),
		userHandler:    userHandler,
		productHandler: productHandler,
		cartHandler:    cartHandler,
		orderHandler:   orderHandler,
	}

	r.setupMiddleware()
	r.setupRoutes()
	return r
}

// setupRoutes configures all routes
func (r *Router) setupRoutes() {
	// Health check
	r.engine.GET("/health", r.healthCheck)
	r.engine.GET("/api/v1/health", r.healthCheck)

	// User routes - Public
	r.engine.POST("/api/v1/users/register", r.userHandler.Register)
	r.engine.POST("/api/v1/users/login", r.userHandler.Login)

	// User routes - Authenticated
	r.engine.GET("/api/v1/users/profile", r.withAuth(), r.userHandler.GetProfile)
	r.engine.PUT("/api/v1/users/update", r.withAuth(), r.userHandler.UpdateUser)

	// User routes - Admin only
	r.engine.GET("/api/v1/users/search", r.withAuth(), r.withRole("admin"), r.userHandler.SearchUsers)
	r.engine.GET("/api/v1/users/by-id", r.withAuth(), r.withRole("admin"), r.userHandler.GetUserByID)
	r.engine.DELETE("/api/v1/users/delete", r.withAuth(), r.withRole("admin"), r.userHandler.DeleteUser)

	// Address routes - Authenticated
	r.engine.POST("/api/v1/addresses/create", r.withAuth(), r.userHandler.CreateAddress)
	r.engine.GET("/api/v1/addresses/list", r.withAuth(), r.userHandler.ListAddresses)
	r.engine.PUT("/api/v1/addresses/update", r.withAuth(), r.userHandler.UpdateAddress)
	r.engine.DELETE("/api/v1/addresses/delete", r.withAuth(), r.userHandler.DeleteAddress)

	// Product routes - Public
	r.engine.GET("/api/v1/products", gin.WrapF(r.productHandler.ListProducts))
	r.engine.GET("/api/v1/products/by-id", gin.WrapF(r.productHandler.GetProductByID))

	// Product routes - Admin only
	r.engine.POST("/api/v1/products/create", r.withAuth(), r.withRole("admin"), gin.WrapF(r.productHandler.CreateProduct))
	r.engine.PUT("/api/v1/products/update", r.withAuth(), r.withRole("admin"), gin.WrapF(r.productHandler.UpdateProduct))
	r.engine.DELETE("/api/v1/products/delete", r.withAuth(), r.withRole("admin"), gin.WrapF(r.productHandler.DeleteProduct))

	// Category routes - Public
	r.engine.GET("/api/v1/categories", gin.WrapF(r.productHandler.ListCategories))
	r.engine.GET("/api/v1/categories/by-id", gin.WrapF(r.productHandler.GetCategoryByID))

	// Category routes - Admin only
	r.engine.POST("/api/v1/categories/create", r.withAuth(), r.withRole("admin"), gin.WrapF(r.productHandler.CreateCategory))
	r.engine.PUT("/api/v1/categories/update", r.withAuth(), r.withRole("admin"), gin.WrapF(r.productHandler.UpdateCategory))
	r.engine.DELETE("/api/v1/categories/delete", r.withAuth(), r.withRole("admin"), gin.WrapF(r.productHandler.DeleteCategory))

	// Cart routes - Authenticated
	r.engine.GET("/api/v1/cart", r.withAuth(), gin.WrapF(r.cartHandler.GetCart))
	r.engine.POST("/api/v1/cart/items/add", r.withAuth(), gin.WrapF(r.cartHandler.AddItem))
	r.engine.PUT("/api/v1/cart/items/update", r.withAuth(), gin.WrapF(r.cartHandler.UpdateItem))
	r.engine.DELETE("/api/v1/cart/items/remove", r.withAuth(), gin.WrapF(r.cartHandler.RemoveItem))
	r.engine.DELETE("/api/v1/cart/clear", r.withAuth(), gin.WrapF(r.cartHandler.ClearCart))

	// Order routes - Authenticated
	r.engine.POST("/api/v1/orders/create", r.withAuth(), gin.WrapF(r.orderHandler.CreateOrder))
	r.engine.GET("/api/v1/orders", r.withAuth(), gin.WrapF(r.orderHandler.ListOrders))
	r.engine.GET("/api/v1/orders/by-id", r.withAuth(), gin.WrapF(r.orderHandler.GetOrderByID))
	r.engine.POST("/api/v1/orders/items/add", r.withAuth(), gin.WrapF(r.orderHandler.AddOrderItem))
	r.engine.DELETE("/api/v1/orders/items/remove", r.withAuth(), gin.WrapF(r.orderHandler.RemoveOrderItem))

	// Order routes - Admin only
	r.engine.PATCH("/api/v1/orders/status", r.withAuth(), r.withRole("admin"), gin.WrapF(r.orderHandler.UpdateOrderStatus))
}

// Handler returns the configured HTTP handler with all middlewares
func (r *Router) Handler() http.Handler {
	return r.engine
}

// Engine exposes the gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) setupMiddleware() {
	r.engine.Use(middleware.CORS(r.cfg.AllowedOrigins, r.cfg.AllowedMethods, r.cfg.AllowedHeaders))
	r.engine.Use(middleware.Recovery())
	r.engine.Use(middleware.RequestID())
	r.engine.Use(middleware.Logger())
	r.engine.Use(middleware.Cancellation())
	r.engine.Use(middleware.Timeout(r.cfg.RequestTimeout))
	r.engine.Use(middleware.NewRateLimiter(r.cfg.RateLimitRequests, r.cfg.RateLimitWindow).Middleware())
}

func (r *Router) withAuth() gin.HandlerFunc {
	return middleware.AuthMiddleware(r.jwtManager)
}

func (r *Router) withRole(roles ...string) gin.HandlerFunc {
	return middleware.RequireRole(roles...)
}

// healthCheck endpoint
func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "api-gateway"})
}
