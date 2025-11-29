package service

import (
	"context"
	"fmt"

	"github.com/a-h-pooladvand/microservices/internal/core/domain"
	"github.com/a-h-pooladvand/microservices/internal/core/port"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Ensure UserService implements port.UserService.
var _ port.UserService = (*UserService)(nil)

// UserService implements the user business logic.
type UserService struct {
	repo   port.UserRepository
	logger *zap.Logger
}

// NewUserService creates a new UserService instance.
func NewUserService(repo port.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context, name, surname string) (*domain.User, error) {
	s.logger.Debug("creating user", zap.String("name", name), zap.String("surname", surname))

	user, err := domain.NewUser(name, surname)
	if err != nil {
		s.logger.Error("failed to create user domain entity", zap.Error(err))
		return nil, fmt.Errorf("create user: %w", err)
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to persist user", zap.Error(err), zap.String("name", name))
		return nil, fmt.Errorf("create user: %w", err)
	}

	s.logger.Info("user created successfully", zap.String("id", createdUser.ID.String()))
	return createdUser, nil
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	s.logger.Debug("getting user", zap.String("id", id.String()))

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

// ListUsers retrieves all users with pagination.
func (s *UserService) ListUsers(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	s.logger.Debug("listing users", zap.Int("offset", offset), zap.Int("limit", limit))

	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.repo.GetAll(ctx, offset, limit)
	if err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil, fmt.Errorf("list users: %w", err)
	}

	return users, nil
}

// UpdateUser updates an existing user.
func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, name, surname string) (*domain.User, error) {
	s.logger.Debug("updating user", zap.String("id", id.String()))

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get user for update", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("update user: %w", err)
	}

	if name != "" {
		user.Name = name
	}
	if surname != "" {
		user.Surname = surname
	}

	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		s.logger.Error("failed to update user", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("update user: %w", err)
	}

	s.logger.Info("user updated successfully", zap.String("id", id.String()))
	return updatedUser, nil
}

// DeleteUser deletes a user by ID.
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	s.logger.Debug("deleting user", zap.String("id", id.String()))

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete user", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("delete user: %w", err)
	}

	s.logger.Info("user deleted successfully", zap.String("id", id.String()))
	return nil
}
