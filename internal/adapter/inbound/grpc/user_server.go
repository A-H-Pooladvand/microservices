package grpc

import (
	"context"
	"errors"

	pb "github.com/a-h-pooladvand/microservices/api/user/v1"
	"github.com/a-h-pooladvand/microservices/internal/core/domain"
	"github.com/a-h-pooladvand/microservices/internal/core/port"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServiceServer implements the gRPC user service.
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	userService port.UserService
}

// NewUserServiceServer creates a new UserServiceServer instance.
func NewUserServiceServer(userService port.UserService) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
	}
}

// CreateUser creates a new user.
func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	if req.GetName() == "" || req.GetSurname() == "" {
		return nil, status.Error(codes.InvalidArgument, "name and surname are required")
	}

	user, err := s.userService.CreateUser(ctx, req.GetName(), req.GetSurname())
	if err != nil {
		if isValidationError(err) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return toProtoUser(user), nil
}

// GetUser retrieves a user by ID.
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	user, err := s.userService.GetUser(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return toProtoUser(user), nil
}

// ListUsers retrieves all users with pagination.
func (s *UserServiceServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	offset := int(req.GetOffset())
	limit := int(req.GetLimit())

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	users, err := s.userService.ListUsers(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	protoUsers := make([]*pb.UserResponse, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(user)
	}

	return &pb.ListUsersResponse{
		Users: protoUsers,
	}, nil
}

// UpdateUser updates an existing user.
func (s *UserServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	user, err := s.userService.UpdateUser(ctx, id, req.GetName(), req.GetSurname())
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if isValidationError(err) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return toProtoUser(user), nil
}

// DeleteUser deletes a user by ID.
func (s *UserServiceServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	err = s.userService.DeleteUser(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: "user deleted successfully",
	}, nil
}

// toProtoUser converts a domain user to a protobuf user response.
func toProtoUser(user *domain.User) *pb.UserResponse {
	return &pb.UserResponse{
		Id:        user.ID.String(),
		Name:      user.Name,
		Surname:   user.Surname,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// isNotFoundError checks if the error is a not found error.
func isNotFoundError(err error) bool {
	return errors.Is(err, domain.ErrUserNotFound)
}

// isValidationError checks if the error is a validation error.
func isValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidUser)
}
