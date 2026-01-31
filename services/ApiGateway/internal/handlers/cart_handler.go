package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kareemhamed001/e-commerce/pkg/logger"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/internal/middleware"
	cartpb "github.com/kareemhamed001/e-commerce/shared/proto/v1/cart"
)

// CartHandler handles cart-related HTTP requests
type CartHandler struct {
	cartClient cartpb.CartServiceClient
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartClient cartpb.CartServiceClient) *CartHandler {
	return &CartHandler{
		cartClient: cartClient,
	}
}

// GetCart godoc
// @Summary Get user cart
// @Description Get the current user's cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CartResponse
// @Router /api/v1/cart [get]
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.cartClient.GetCart(r.Context(), &cartpb.GetCartRequest{
		UserId: int64(userID),
	})

	if err != nil {
		logger.Errorf("failed to get cart: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// AddItem godoc
// @Summary Add item to cart
// @Description Add a product to the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddItemRequest true "Item details"
// @Success 200 {object} CartResponse
// @Router /api/v1/cart/items [post]
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		ProductID int64 `json:"product_id"`
		Quantity  int32 `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.cartClient.AddItem(r.Context(), &cartpb.AddItemRequest{
		UserId:    int64(userID),
		ProductId: req.ProductID,
		Quantity:  req.Quantity,
	})

	if err != nil {
		logger.Errorf("failed to add item to cart: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// UpdateItem godoc
// @Summary Update cart item
// @Description Update the quantity of a cart item
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateItemRequest true "Item update details"
// @Success 200 {object} CartResponse
// @Router /api/v1/cart/items [put]
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		ProductID int64 `json:"product_id"`
		Quantity  int32 `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.cartClient.UpdateItem(r.Context(), &cartpb.UpdateItemRequest{
		UserId:    int64(userID),
		ProductId: req.ProductID,
		Quantity:  req.Quantity,
	})

	if err != nil {
		logger.Errorf("failed to update cart item: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// RemoveItem godoc
// @Summary Remove item from cart
// @Description Remove a product from the user's cart
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RemoveItemRequest true "Product ID"
// @Success 200 {object} CartResponse
// @Router /api/v1/cart/items [delete]
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		ProductID int64 `json:"product_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.cartClient.RemoveItem(r.Context(), &cartpb.RemoveItemRequest{
		UserId:    int64(userID),
		ProductId: req.ProductID,
	})

	if err != nil {
		logger.Errorf("failed to remove item from cart: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ClearCart godoc
// @Summary Clear cart
// @Description Remove all items from the user's cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ClearCartResponse
// @Router /api/v1/cart [delete]
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.cartClient.ClearCart(r.Context(), &cartpb.ClearCartRequest{
		UserId: int64(userID),
	})

	if err != nil {
		logger.Errorf("failed to clear cart: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
