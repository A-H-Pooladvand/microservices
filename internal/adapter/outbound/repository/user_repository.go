package repository

import (
	"context"
	"errors"
	"time"

	"github.com/a-h-pooladvand/microservices/internal/core/domain"
	"github.com/a-h-pooladvand/microservices/internal/core/port"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ensure PostgresUserRepository implements port.UserRepository.
var _ port.UserRepository = (*PostgresUserRepository)(nil)

// UserModel represents the database model for user.
type UserModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `gorm:"not null"`
	Surname   string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for user model.
func (UserModel) TableName() string {
	return "users"
}

// toDomain converts database model to domain entity.
func (m *UserModel) toDomain() *domain.User {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	return &domain.User{
		ID:        m.ID,
		Name:      m.Name,
		Surname:   m.Surname,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

// fromDomain converts domain entity to database model.
func fromDomain(u *domain.User) *UserModel {
	return &UserModel{
		ID:        u.ID,
		Name:      u.Name,
		Surname:   u.Surname,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// PostgresUserRepository implements the user repository using PostgreSQL.
type PostgresUserRepository struct {
	db *gorm.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository instance.
func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// Create creates a new user in the database.
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	model := fromDomain(user)

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("failed to create user: no rows affected")
	}

	return model.toDomain(), nil
}

// GetByID retrieves a user by ID.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var model UserModel

	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, result.Error
	}

	return model.toDomain(), nil
}

// GetAll retrieves all users with pagination.
func (r *PostgresUserRepository) GetAll(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	var models []UserModel

	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = model.toDomain()
	}

	return users, nil
}

// Update updates an existing user.
func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	model := fromDomain(user)
	model.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, domain.ErrUserNotFound
	}

	return model.toDomain(), nil
}

// Delete soft deletes a user by ID.
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&UserModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Migrate runs database migrations for the user model.
func (r *PostgresUserRepository) Migrate() error {
	return r.db.AutoMigrate(&UserModel{})
}
