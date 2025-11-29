package port

import (
	"context"

	"github.com/a-h-pooladvand/microservices/internal/core/domain"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user persistence operations.
// This is a driven port (secondary/outbound port).
type UserRepository interface {
	// Create creates a new user in the repository.
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	// GetByID retrieves a user by ID.
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	// GetAll retrieves all users with pagination.
	GetAll(ctx context.Context, offset, limit int) ([]*domain.User, error)
	// Update updates an existing user.
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	// Delete deletes a user by ID.
	Delete(ctx context.Context, id uuid.UUID) error
}
