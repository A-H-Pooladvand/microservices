package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/a-h-pooladvand/microservices/internal/core/domain"
	"github.com/a-h-pooladvand/microservices/internal/core/service"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// MockUserRepository is a mock implementation of port.UserRepository.
type MockUserRepository struct {
	users    map[uuid.UUID]*domain.User
	createFn func(ctx context.Context, user *domain.User) (*domain.User, error)
	getByIDFn func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	getAllFn  func(ctx context.Context, offset, limit int) ([]*domain.User, error)
	updateFn  func(ctx context.Context, user *domain.User) (*domain.User, error)
	deleteFn  func(ctx context.Context, id uuid.UUID) error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uuid.UUID]*domain.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	user, ok := m.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) GetAll(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, offset, limit)
	}
	var users []*domain.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	if _, ok := m.users[user.ID]; !ok {
		return nil, domain.ErrUserNotFound
	}
	user.UpdatedAt = time.Now()
	m.users[user.ID] = user
	return user, nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	if _, ok := m.users[id]; !ok {
		return domain.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func TestUserService_CreateUser(t *testing.T) {
	logger := zap.NewNop()
	
	tests := []struct {
		name      string
		userName  string
		surname   string
		setupMock func(*MockUserRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "successful user creation",
			userName: "John",
			surname:  "Doe",
			wantErr:  false,
		},
		{
			name:     "empty name should fail",
			userName: "",
			surname:  "Doe",
			wantErr:  true,
			errMsg:   "invalid user data",
		},
		{
			name:     "empty surname should fail",
			userName: "John",
			surname:  "",
			wantErr:  true,
			errMsg:   "invalid user data",
		},
		{
			name:     "repository error should propagate",
			userName: "John",
			surname:  "Doe",
			setupMock: func(repo *MockUserRepository) {
				repo.createFn = func(ctx context.Context, user *domain.User) (*domain.User, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			
			svc := service.NewUserService(repo, logger)
			
			user, err := svc.CreateUser(context.Background(), tt.userName, tt.surname)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateUser() expected error but got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("CreateUser() unexpected error: %v", err)
				return
			}
			
			if user == nil {
				t.Error("CreateUser() returned nil user")
				return
			}
			
			if user.Name != tt.userName {
				t.Errorf("CreateUser() name = %v, want %v", user.Name, tt.userName)
			}
			
			if user.Surname != tt.surname {
				t.Errorf("CreateUser() surname = %v, want %v", user.Surname, tt.surname)
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	logger := zap.NewNop()
	
	existingUser := &domain.User{
		ID:        uuid.New(),
		Name:      "Jane",
		Surname:   "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	tests := []struct {
		name      string
		userID    uuid.UUID
		setupMock func(*MockUserRepository)
		wantErr   bool
		wantUser  *domain.User
	}{
		{
			name:   "existing user",
			userID: existingUser.ID,
			setupMock: func(repo *MockUserRepository) {
				repo.users[existingUser.ID] = existingUser
			},
			wantErr:  false,
			wantUser: existingUser,
		},
		{
			name:    "non-existing user",
			userID:  uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			
			svc := service.NewUserService(repo, logger)
			
			user, err := svc.GetUser(context.Background(), tt.userID)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetUser() expected error but got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("GetUser() unexpected error: %v", err)
				return
			}
			
			if user.ID != tt.wantUser.ID {
				t.Errorf("GetUser() id = %v, want %v", user.ID, tt.wantUser.ID)
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	logger := zap.NewNop()
	
	tests := []struct {
		name      string
		offset    int
		limit     int
		setupMock func(*MockUserRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:   "list with default pagination",
			offset: 0,
			limit:  0,
			setupMock: func(repo *MockUserRepository) {
				for i := 0; i < 5; i++ {
					id := uuid.New()
					repo.users[id] = &domain.User{
						ID:        id,
						Name:      "User",
						Surname:   "Test",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
				}
			},
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "empty list",
			offset:    0,
			limit:     10,
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:   "negative offset should be normalized",
			offset: -10,
			limit:  10,
			setupMock: func(repo *MockUserRepository) {
				id := uuid.New()
				repo.users[id] = &domain.User{
					ID: id, Name: "User", Surname: "Test",
					CreatedAt: time.Now(), UpdatedAt: time.Now(),
				}
			},
			wantCount: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			
			svc := service.NewUserService(repo, logger)
			
			users, err := svc.ListUsers(context.Background(), tt.offset, tt.limit)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ListUsers() expected error but got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("ListUsers() unexpected error: %v", err)
				return
			}
			
			if len(users) != tt.wantCount {
				t.Errorf("ListUsers() count = %v, want %v", len(users), tt.wantCount)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	logger := zap.NewNop()
	
	existingUser := &domain.User{
		ID:        uuid.New(),
		Name:      "Original",
		Surname:   "Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	tests := []struct {
		name        string
		userID      uuid.UUID
		newName     string
		newSurname  string
		setupMock   func(*MockUserRepository)
		wantErr     bool
		wantName    string
		wantSurname string
	}{
		{
			name:       "update name only",
			userID:     existingUser.ID,
			newName:    "Updated",
			newSurname: "",
			setupMock: func(repo *MockUserRepository) {
				repo.users[existingUser.ID] = &domain.User{
					ID:        existingUser.ID,
					Name:      existingUser.Name,
					Surname:   existingUser.Surname,
					CreatedAt: existingUser.CreatedAt,
					UpdatedAt: existingUser.UpdatedAt,
				}
			},
			wantErr:     false,
			wantName:    "Updated",
			wantSurname: "Name",
		},
		{
			name:       "update surname only",
			userID:     existingUser.ID,
			newName:    "",
			newSurname: "NewSurname",
			setupMock: func(repo *MockUserRepository) {
				repo.users[existingUser.ID] = &domain.User{
					ID:        existingUser.ID,
					Name:      existingUser.Name,
					Surname:   existingUser.Surname,
					CreatedAt: existingUser.CreatedAt,
					UpdatedAt: existingUser.UpdatedAt,
				}
			},
			wantErr:     false,
			wantName:    "Original",
			wantSurname: "NewSurname",
		},
		{
			name:       "user not found",
			userID:     uuid.New(),
			newName:    "Updated",
			newSurname: "Name",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			
			svc := service.NewUserService(repo, logger)
			
			user, err := svc.UpdateUser(context.Background(), tt.userID, tt.newName, tt.newSurname)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateUser() expected error but got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("UpdateUser() unexpected error: %v", err)
				return
			}
			
			if user.Name != tt.wantName {
				t.Errorf("UpdateUser() name = %v, want %v", user.Name, tt.wantName)
			}
			
			if user.Surname != tt.wantSurname {
				t.Errorf("UpdateUser() surname = %v, want %v", user.Surname, tt.wantSurname)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	logger := zap.NewNop()
	
	existingUser := &domain.User{
		ID:        uuid.New(),
		Name:      "ToDelete",
		Surname:   "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	tests := []struct {
		name      string
		userID    uuid.UUID
		setupMock func(*MockUserRepository)
		wantErr   bool
	}{
		{
			name:   "delete existing user",
			userID: existingUser.ID,
			setupMock: func(repo *MockUserRepository) {
				repo.users[existingUser.ID] = existingUser
			},
			wantErr: false,
		},
		{
			name:    "delete non-existing user",
			userID:  uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockUserRepository()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}
			
			svc := service.NewUserService(repo, logger)
			
			err := svc.DeleteUser(context.Background(), tt.userID)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteUser() expected error but got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("DeleteUser() unexpected error: %v", err)
			}
			
			// Verify user is actually deleted
			_, err = svc.GetUser(context.Background(), tt.userID)
			if err == nil {
				t.Error("DeleteUser() user should be deleted but still exists")
			}
		})
	}
}

func TestDomain_NewUser(t *testing.T) {
	tests := []struct {
		name    string
		userName string
		surname string
		wantErr bool
	}{
		{
			name:    "valid user",
			userName: "John",
			surname: "Doe",
			wantErr: false,
		},
		{
			name:    "empty name",
			userName: "",
			surname: "Doe",
			wantErr: true,
		},
		{
			name:    "empty surname",
			userName: "John",
			surname: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := domain.NewUser(tt.userName, tt.surname)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error but got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("NewUser() unexpected error: %v", err)
				return
			}
			
			if user.Name != tt.userName {
				t.Errorf("NewUser() name = %v, want %v", user.Name, tt.userName)
			}
			
			if user.Surname != tt.surname {
				t.Errorf("NewUser() surname = %v, want %v", user.Surname, tt.surname)
			}
			
			if user.ID == uuid.Nil {
				t.Error("NewUser() ID should not be nil")
			}
		})
	}
}

func TestDomain_User_FullName(t *testing.T) {
	user := &domain.User{
		ID:      uuid.New(),
		Name:    "John",
		Surname: "Doe",
	}
	
	expected := "John Doe"
	if got := user.FullName(); got != expected {
		t.Errorf("FullName() = %v, want %v", got, expected)
	}
}
