package port

import (
	"context"

	"github.com/a-h-pooladvand/microservices/internal/core/domain"
	"github.com/google/uuid"
)

// UserService defines the interface for user business logic.
// This is a driver port (primary/inbound port).
type UserService interface {
	// CreateUser creates a new user.
	CreateUser(ctx context.Context, name, surname string) (*domain.User, error)
	// GetUser retrieves a user by ID.
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	// ListUsers retrieves all users with pagination.
	ListUsers(ctx context.Context, offset, limit int) ([]*domain.User, error)
	// UpdateUser updates an existing user.
	UpdateUser(ctx context.Context, id uuid.UUID, name, surname string) (*domain.User, error)
	// DeleteUser deletes a user by ID.
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
