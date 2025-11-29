package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrUserNotFound is returned when a user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidUser is returned when user data is invalid.
	ErrInvalidUser = errors.New("invalid user data")
	// ErrUserAlreadyExists is returned when a user already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
)

// User represents the user domain entity.
type User struct {
	ID        uuid.UUID
	Name      string
	Surname   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewUser creates a new user domain entity.
func NewUser(name, surname string) (*User, error) {
	if name == "" || surname == "" {
		return nil, ErrInvalidUser
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Surname:   surname,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Validate validates the user entity.
func (u *User) Validate() error {
	if u.Name == "" || u.Surname == "" {
		return ErrInvalidUser
	}
	return nil
}

// FullName returns the full name of the user.
func (u *User) FullName() string {
	return u.Name + " " + u.Surname
}
