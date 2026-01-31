package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kareemhamed001/e-commerce/pkg/logger"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/internal/middleware"
	orderpb "github.com/kareemhamed001/e-commerce/shared/proto/v1/order"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderClient orderpb.OrderServiceClient
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderClient orderpb.OrderServiceClient) *OrderHandler {
	return &OrderHandler{
		orderClient: orderClient,
	}
}

// CreateOrder godoc
// @Summary Create order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrderRequest true "Order details"
// @Success 201 {object} CreateOrderResponse
// @Router /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		ShippingCost         float32 `json:"shipping_cost"`
		ShippingDurationDays int32   `json:"shipping_duration_days"`
		Discount             float32 `json:"discount"`
		Items                []struct {
			ProductID int64 `json:"product_id"`
			Quantity  int32 `json:"quantity"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	items := make([]*orderpb.OrderItemInput, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, &orderpb.OrderItemInput{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	resp, err := h.orderClient.CreateOrder(r.Context(), &orderpb.CreateOrderRequest{
		UserId:               int64(userID),
		ShippingCost:         req.ShippingCost,
		ShippingDurationDays: req.ShippingDurationDays,
		Discount:             req.Discount,
		Items:                items,
	})
	if err != nil {
		logger.Errorf("failed to create order: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Get order details by ID
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id query int true "Order ID"
// @Success 200 {object} GetOrderByIDResponse
// @Router /api/v1/orders/by-id [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "missing order ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid order ID")
		return
	}

	resp, err := h.orderClient.GetOrderByID(r.Context(), &orderpb.GetOrderByIDRequest{
		Id: id,
	})
	if err != nil {
		logger.Errorf("failed to get order: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ListOrders godoc
// @Summary List orders
// @Description List orders with pagination
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param user_id query int false "Filter by user ID (admin only)"
// @Success 200 {object} ListOrdersResponse
// @Router /api/v1/orders [get]
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	userIDParam := r.URL.Query().Get("user_id")
	var userIDFilter int64
	if userIDParam != "" {
		userIDFilter, _ = strconv.ParseInt(userIDParam, 10, 64)
	} else {
		userID, ok := middleware.GetUserID(r.Context())
		if ok {
			userIDFilter = int64(userID)
		}
	}

	resp, err := h.orderClient.ListOrders(r.Context(), &orderpb.ListOrdersRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
		UserId:  userIDFilter,
	})
	if err != nil {
		logger.Errorf("failed to list orders: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// AddOrderItem godoc
// @Summary Add item to order
// @Description Add a new item to an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddOrderItemRequest true "Order item details"
// @Success 200 {object} AddOrderItemResponse
// @Router /api/v1/orders/items/add [post]
func (h *OrderHandler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	var req orderpb.AddOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.orderClient.AddOrderItem(r.Context(), &req)
	if err != nil {
		logger.Errorf("failed to add order item: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// RemoveOrderItem godoc
// @Summary Remove item from order
// @Description Remove an item from an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RemoveOrderItemRequest true "Order item ID"
// @Success 200 {object} RemoveOrderItemResponse
// @Router /api/v1/orders/items/remove [delete]
func (h *OrderHandler) RemoveOrderItem(w http.ResponseWriter, r *http.Request) {
	var req orderpb.RemoveOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.orderClient.RemoveOrderItem(r.Context(), &req)
	if err != nil {
		logger.Errorf("failed to remove order item: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update the status of an order (admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateOrderStatusRequest true "Status update details"
// @Success 200 {object} UpdateOrderStatusResponse
// @Router /api/v1/orders/status [patch]
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	var req orderpb.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.orderClient.UpdateOrderStatus(r.Context(), &req)
	if err != nil {
		logger.Errorf("failed to update order status: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
