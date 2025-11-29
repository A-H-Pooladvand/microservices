package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/a-h-pooladvand/microservices/internal/core/domain"
	"github.com/a-h-pooladvand/microservices/internal/core/port"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	userService port.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userService port.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest represents the request body for creating a user.
type CreateUserRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=50"`
	Surname string `json:"surname" validate:"required,min=2,max=50"`
}

// UpdateUserRequest represents the request body for updating a user.
type UpdateUserRequest struct {
	Name    string `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	Surname string `json:"surname,omitempty" validate:"omitempty,min=2,max=50"`
}

// UserResponse represents the response for user operations.
type UserResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Surname   string  `json:"surname"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

// APIResponse represents a standard API response.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// toUserResponse converts a domain user to a response.
func toUserResponse(user *domain.User) *UserResponse {
	var deletedAt *string
	if user.DeletedAt != nil {
		t := user.DeletedAt.Format("2006-01-02T15:04:05Z07:00")
		deletedAt = &t
	}
	return &UserResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Surname:   user.Surname,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		DeletedAt: deletedAt,
	}
}

// Create handles the creation of a new user.
//
//	@Summary		Create a new user
//	@Description	Creates a new user in the system
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateUserRequest	true	"User creation request"
//	@Success		201		{object}	APIResponse{data=UserResponse}
//	@Failure		400		{object}	APIResponse
//	@Failure		422		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/users [post]
func (h *UserHandler) Create(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	user, err := h.userService.CreateUser(c.Request().Context(), req.Name, req.Surname)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUser) {
			return c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "invalid user data",
			})
		}
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    toUserResponse(user),
		Message: "user created successfully",
	})
}

// Get handles retrieving a user by ID.
//
//	@Summary		Get a user
//	@Description	Retrieves a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	APIResponse{data=UserResponse}
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/users/{id} [get]
func (h *UserHandler) Get(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "invalid user ID format",
		})
	}

	user, err := h.userService.GetUser(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Error:   "user not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to retrieve user",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    toUserResponse(user),
	})
}

// List handles listing all users with pagination.
//
//	@Summary		List users
//	@Description	Retrieves all users with pagination
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			offset	query		int	false	"Offset for pagination"	default(0)
//	@Param			limit	query		int	false	"Limit for pagination"	default(10)
//	@Success		200		{object}	APIResponse{data=[]UserResponse}
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/users [get]
func (h *UserHandler) List(c echo.Context) error {
	offset := 0
	limit := 10

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := parsePositiveInt(o); err == nil {
			offset = parsed
		}
	}
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := parsePositiveInt(l); err == nil {
			limit = parsed
		}
	}

	users, err := h.userService.ListUsers(c.Request().Context(), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to list users",
		})
	}

	responses := make([]*UserResponse, len(users))
	for i, user := range users {
		responses[i] = toUserResponse(user)
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    responses,
	})
}

// Update handles updating a user.
//
//	@Summary		Update a user
//	@Description	Updates an existing user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"User ID"
//	@Param			request	body		UpdateUserRequest	true	"User update request"
//	@Success		200		{object}	APIResponse{data=UserResponse}
//	@Failure		400		{object}	APIResponse
//	@Failure		404		{object}	APIResponse
//	@Failure		422		{object}	APIResponse
//	@Failure		500		{object}	APIResponse
//	@Router			/api/v1/users/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "invalid user ID format",
		})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "invalid request body",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), id, req.Name, req.Surname)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Error:   "user not found",
			})
		}
		if errors.Is(err, domain.ErrInvalidUser) {
			return c.JSON(http.StatusBadRequest, APIResponse{
				Success: false,
				Error:   "invalid user data",
			})
		}
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to update user",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    toUserResponse(user),
		Message: "user updated successfully",
	})
}

// Delete handles deleting a user.
//
//	@Summary		Delete a user
//	@Description	Deletes a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	APIResponse
//	@Failure		400	{object}	APIResponse
//	@Failure		404	{object}	APIResponse
//	@Failure		500	{object}	APIResponse
//	@Router			/api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "invalid user ID format",
		})
	}

	err = h.userService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, APIResponse{
				Success: false,
				Error:   "user not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to delete user",
		})
	}

	return c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "user deleted successfully",
	})
}

// parsePositiveInt parses a string to a positive integer.
func parsePositiveInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0, errors.New("invalid positive integer")
	}
	return n, nil
}
