package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kareemhamed001/e-commerce/pkg/logger"
	productpb "github.com/kareemhamed001/e-commerce/shared/proto/v1/product"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productClient productpb.ProductServiceClient
}

// NewProductHandler creates a new product handler
func NewProductHandler(productClient productpb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{
		productClient: productClient,
	}
}

// CreateProduct godoc
// @Summary Create product
// @Description Create a new product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProductRequest true "Product details"
// @Success 201 {object} CreateProductResponse
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req productpb.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.productClient.CreateProduct(r.Context(), &req)
	if err != nil {
		logger.Errorf("failed to create product: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} GetProductByIDResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "missing product ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	resp, err := h.productClient.GetProductByID(r.Context(), &productpb.GetProductByIDRequest{
		Id: id,
	})

	if err != nil {
		logger.Errorf("failed to get product: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ListProducts godoc
// @Summary List products
// @Description List all products with pagination
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} ListProductsResponse
// @Router /api/v1/products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	resp, err := h.productClient.ListProducts(r.Context(), &productpb.ListProductsRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
	})

	if err != nil {
		logger.Errorf("failed to list products: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product details (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProductRequest true "Product update details"
// @Success 200 {object} UpdateProductResponse
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var req productpb.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.productClient.UpdateProduct(r.Context(), &req)
	if err != nil {
		logger.Errorf("failed to update product: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product (admin only)
// @Tags products
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} DeleteProductResponse
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "missing product ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	resp, err := h.productClient.DeleteProduct(r.Context(), &productpb.DeleteProductRequest{
		Id: id,
	})

	if err != nil {
		logger.Errorf("failed to delete product: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// Category handlers

// CreateCategory godoc
// @Summary Create category
// @Description Create a new category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCategoryRequest true "Category details"
// @Success 201 {object} CreateCategoryResponse
// @Router /api/v1/categories [post]
func (h *ProductHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req productpb.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.productClient.CreateCategory(r.Context(), &req)
	if err != nil {
		logger.Errorf("failed to create category: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Get category details by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} GetCategoryByIDResponse
// @Router /api/v1/categories/{id} [get]
func (h *ProductHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "missing category ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	resp, err := h.productClient.GetCategoryByID(r.Context(), &productpb.GetCategoryByIDRequest{
		Id: id,
	})

	if err != nil {
		logger.Errorf("failed to get category: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ListCategories godoc
// @Summary List categories
// @Description List all categories with pagination
// @Tags categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} ListCategoriesResponse
// @Router /api/v1/categories [get]
func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	resp, err := h.productClient.ListCategories(r.Context(), &productpb.ListCategoriesRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
	})

	if err != nil {
		logger.Errorf("failed to list categories: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update category details (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateCategoryRequest true "Category update details"
// @Success 200 {object} UpdateCategoryResponse
// @Router /api/v1/categories/{id} [put]
func (h *ProductHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var req productpb.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.productClient.UpdateCategory(r.Context(), &req)
	if err != nil {
		logger.Errorf("failed to update category: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete a category (admin only)
// @Tags categories
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} DeleteCategoryResponse
// @Router /api/v1/categories/{id} [delete]
func (h *ProductHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, http.StatusBadRequest, "missing category ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	resp, err := h.productClient.DeleteCategory(r.Context(), &productpb.DeleteCategoryRequest{
		Id: id,
	})

	if err != nil {
		logger.Errorf("failed to delete category: %v", err)
		writeJSONErrorFromGRPC(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
