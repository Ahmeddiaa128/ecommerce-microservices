package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kareemhamed001/e-commerce/pkg/logger"
	"github.com/kareemhamed001/e-commerce/services/ApiGateway/internal/middleware"
	userpb "github.com/kareemhamed001/e-commerce/shared/proto/v1/user"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userClient userpb.UserServiceClient
}

// NewUserHandler creates a new user handler
func NewUserHandler(userClient userpb.UserServiceClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User registration details"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/users/register [post]

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, "invalid request body")
		return
	}

	// Default role to "customer" if not specified
	if req.Role == "" {
		req.Role = "customer"
	}

	resp, err := h.userClient.CreateUser(c.Request.Context(), &userpb.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})

	if err != nil {
		logger.Errorf("failed to create user: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.userClient.Login(c.Request.Context(), &userpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		logger.Errorf("login failed: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} User
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c.Request.Context())
	if !ok {
		writeJSONError(c.Writer, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.userClient.GetUserByID(c.Request.Context(), &userpb.GetUserByIDRequest{
		Id: int32(userID),
	})

	if err != nil {
		logger.Errorf("failed to get user: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user details by ID (admin or self)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		writeJSONError(c.Writer, http.StatusBadRequest, "missing user ID")
		return
	}

	parsedID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, "invalid user ID")
		return
	}
	id := parsedID

	resp, err := h.userClient.GetUserByID(c.Request.Context(), &userpb.GetUserByIDRequest{
		Id: int32(id),
	})

	if err != nil {
		logger.Errorf("failed to get user: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchUsers godoc
// @Summary Search users
// @Description Search users with pagination (admin only)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} SearchUsersResponse
// @Router /api/v1/users [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.Query("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	query := c.Query("query")

	resp, err := h.userClient.SearchUsers(c.Request.Context(), &userpb.SearchUsersRequest{
		Query:      query,
		PageNumber: int32(page),
		PageSize:   int32(perPage),
	})

	if err != nil {
		logger.Errorf("failed to search users: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateUserRequest true "User update details"
// @Success 200 {object} User
// @Router /api/v1/users [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c.Request.Context())
	if !ok {
		writeJSONError(c.Writer, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.userClient.UpdateUser(c.Request.Context(), &userpb.UpdateUserRequest{
		Id:    int32(userID),
		Name:  req.Name,
		Email: req.Email,
	})

	if err != nil {
		logger.Errorf("failed to update user: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user account (admin only)
// @Tags users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} DeleteUserResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		writeJSONError(c.Writer, http.StatusBadRequest, "missing user ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, "invalid user ID")
		return
	}

	resp, err := h.userClient.DeleteUser(c.Request.Context(), &userpb.DeleteUserRequest{
		Id: int32(id),
	})

	if err != nil {
		logger.Errorf("failed to delete user: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Address handlers

// CreateAddress godoc
// @Summary Create address
// @Description Create a new address for authenticated user
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAddressRequest true "Address details"
// @Success 201 {object} CreateAddressResponse
// @Router /api/v1/addresses [post]
func (h *UserHandler) CreateAddress(c *gin.Context) {
	userID, ok := middleware.GetUserID(c.Request.Context())
	if !ok {
		writeJSONError(c.Writer, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req userpb.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, "invalid request body")
		return
	}

	req.UserId = int32(userID)

	resp, err := h.userClient.CreateAddress(c.Request.Context(), &req)
	if err != nil {
		logger.Errorf("failed to create address: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// ListAddresses godoc
// @Summary List user addresses
// @Description Get all addresses for authenticated user
// @Tags addresses
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ListAddressesByUserIDResponse
// @Router /api/v1/addresses [get]
func (h *UserHandler) ListAddresses(c *gin.Context) {
	userID, ok := middleware.GetUserID(c.Request.Context())
	if !ok {
		writeJSONError(c.Writer, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.userClient.ListAddressesByUserID(c.Request.Context(), &userpb.ListAddressesByUserIDRequest{
		UserId: int32(userID),
	})

	if err != nil {
		logger.Errorf("failed to list addresses: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateAddress godoc
// @Summary Update address
// @Description Update an existing address
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateAddressRequest true "Address update details"
// @Success 200 {object} UpdateAddressResponse
// @Router /api/v1/addresses/{id} [put]
func (h *UserHandler) UpdateAddress(c *gin.Context) {
	var req userpb.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.userClient.UpdateAddress(c.Request.Context(), &req)
	if err != nil {
		logger.Errorf("failed to update address: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteAddress godoc
// @Summary Delete address
// @Description Delete an address
// @Tags addresses
// @Security BearerAuth
// @Param id path int true "Address ID"
// @Success 200 {object} DeleteAddressResponse
// @Router /api/v1/addresses/{id} [delete]
func (h *UserHandler) DeleteAddress(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		writeJSONError(c.Writer, http.StatusBadRequest, "missing address ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(c.Writer, http.StatusBadRequest, "invalid address ID")
		return
	}
	address, err := h.userClient.GetAddressByID(c.Request.Context(), &userpb.GetAddressByIDRequest{
		Id: int32(id),
	})
	if err != nil {
		logger.Errorf("failed to get address: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusInternalServerError)
		return
	}

	userID, ok := middleware.GetUserID(c.Request.Context())
	if !ok || address.Address.UserId != int32(userID) {
		writeJSONError(c.Writer, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.userClient.DeleteAddress(c.Request.Context(), &userpb.DeleteAddressRequest{
		Id: int32(id),
	})

	if err != nil {
		logger.Errorf("failed to delete address: %v", err)
		writeJSONErrorFromGRPC(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, resp)
}
