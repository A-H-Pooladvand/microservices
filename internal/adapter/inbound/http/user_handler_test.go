package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httphandler "github.com/a-h-pooladvand/microservices/internal/adapter/inbound/http"
	"github.com/a-h-pooladvand/microservices/internal/core/domain"
	"github.com/a-h-pooladvand/microservices/pkg/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// MockUserService is a mock implementation of port.UserService.
type MockUserService struct {
	createUserFn func(ctx context.Context, name, surname string) (*domain.User, error)
	getUserFn    func(ctx context.Context, id uuid.UUID) (*domain.User, error)
	listUsersFn  func(ctx context.Context, offset, limit int) ([]*domain.User, error)
	updateUserFn func(ctx context.Context, id uuid.UUID, name, surname string) (*domain.User, error)
	deleteUserFn func(ctx context.Context, id uuid.UUID) error
}

func (m *MockUserService) CreateUser(ctx context.Context, name, surname string) (*domain.User, error) {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, name, surname)
	}
	return nil, nil
}

func (m *MockUserService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if m.getUserFn != nil {
		return m.getUserFn(ctx, id)
	}
	return nil, nil
}

func (m *MockUserService) ListUsers(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	if m.listUsersFn != nil {
		return m.listUsersFn(ctx, offset, limit)
	}
	return nil, nil
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, name, surname string) (*domain.User, error) {
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, id, name, surname)
	}
	return nil, nil
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if m.deleteUserFn != nil {
		return m.deleteUserFn(ctx, id)
	}
	return nil
}

func setupEcho() *echo.Echo {
	e := echo.New()
	e.Validator = validator.New()
	return e
}

func TestUserHandler_Create(t *testing.T) {
	now := time.Now()
	createdUser := &domain.User{
		ID:        uuid.New(),
		Name:      "John",
		Surname:   "Doe",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockService    *MockUserService
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name: "successful creation",
			requestBody: map[string]interface{}{
				"name":    "John",
				"surname": "Doe",
			},
			mockService: &MockUserService{
				createUserFn: func(ctx context.Context, name, surname string) (*domain.User, error) {
					return createdUser, nil
				},
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if !body["success"].(bool) {
					t.Error("expected success to be true")
				}
				data := body["data"].(map[string]interface{})
				if data["name"] != "John" {
					t.Errorf("expected name John, got %v", data["name"])
				}
			},
		},
		{
			name: "validation error - missing name",
			requestBody: map[string]interface{}{
				"surname": "Doe",
			},
			mockService:    &MockUserService{},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				if body["success"].(bool) {
					t.Error("expected success to be false")
				}
			},
		},
		{
			name: "validation error - name too short",
			requestBody: map[string]interface{}{
				"name":    "J",
				"surname": "Doe",
			},
			mockService:    &MockUserService{},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "service error",
			requestBody: map[string]interface{}{
				"name":    "John",
				"surname": "Doe",
			},
			mockService: &MockUserService{
				createUserFn: func(ctx context.Context, name, surname string) (*domain.User, error) {
					return nil, domain.ErrInternal
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			handler := httphandler.NewUserHandler(tt.mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Create(c)
			if err != nil {
				t.Fatalf("handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.checkResponse != nil {
				var response map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	now := time.Now()
	existingUser := &domain.User{
		ID:        uuid.New(),
		Name:      "Jane",
		Surname:   "Smith",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name           string
		userID         string
		mockService    *MockUserService
		expectedStatus int
	}{
		{
			name:   "successful get",
			userID: existingUser.ID.String(),
			mockService: &MockUserService{
				getUserFn: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
					return existingUser, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid uuid",
			userID:         "invalid-uuid",
			mockService:    &MockUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "user not found",
			userID: uuid.New().String(),
			mockService: &MockUserService{
				getUserFn: func(ctx context.Context, id uuid.UUID) (*domain.User, error) {
					return nil, domain.ErrUserNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			handler := httphandler.NewUserHandler(tt.mockService)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			err := handler.Get(c)
			if err != nil {
				t.Fatalf("handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestUserHandler_List(t *testing.T) {
	now := time.Now()
	users := []*domain.User{
		{ID: uuid.New(), Name: "User1", Surname: "Test", CreatedAt: now, UpdatedAt: now},
		{ID: uuid.New(), Name: "User2", Surname: "Test", CreatedAt: now, UpdatedAt: now},
	}

	tests := []struct {
		name           string
		queryParams    map[string]string
		mockService    *MockUserService
		expectedStatus int
		checkResponse  func(t *testing.T, body map[string]interface{})
	}{
		{
			name:        "list all users",
			queryParams: map[string]string{},
			mockService: &MockUserService{
				listUsersFn: func(ctx context.Context, offset, limit int) ([]*domain.User, error) {
					return users, nil
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				if len(data) != 2 {
					t.Errorf("expected 2 users, got %d", len(data))
				}
			},
		},
		{
			name: "list with pagination",
			queryParams: map[string]string{
				"offset": "0",
				"limit":  "10",
			},
			mockService: &MockUserService{
				listUsersFn: func(ctx context.Context, offset, limit int) ([]*domain.User, error) {
					if offset != 0 || limit != 10 {
						t.Errorf("expected offset=0, limit=10, got offset=%d, limit=%d", offset, limit)
					}
					return users, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "empty list",
			queryParams: map[string]string{},
			mockService: &MockUserService{
				listUsersFn: func(ctx context.Context, offset, limit int) ([]*domain.User, error) {
					return []*domain.User{}, nil
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].([]interface{})
				if len(data) != 0 {
					t.Errorf("expected 0 users, got %d", len(data))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			handler := httphandler.NewUserHandler(tt.mockService)

			url := "/api/v1/users"
			if len(tt.queryParams) > 0 {
				url += "?"
				for k, v := range tt.queryParams {
					url += k + "=" + v + "&"
				}
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.List(c)
			if err != nil {
				t.Fatalf("handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.checkResponse != nil {
				var response map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name           string
		userID         string
		mockService    *MockUserService
		expectedStatus int
	}{
		{
			name:   "successful delete",
			userID: userID.String(),
			mockService: &MockUserService{
				deleteUserFn: func(ctx context.Context, id uuid.UUID) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid uuid",
			userID:         "invalid-uuid",
			mockService:    &MockUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "user not found",
			userID: uuid.New().String(),
			mockService: &MockUserService{
				deleteUserFn: func(ctx context.Context, id uuid.UUID) error {
					return domain.ErrUserNotFound
				},
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := setupEcho()
			handler := httphandler.NewUserHandler(tt.mockService)

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+tt.userID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.userID)

			err := handler.Delete(c)
			if err != nil {
				t.Fatalf("handler returned error: %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
