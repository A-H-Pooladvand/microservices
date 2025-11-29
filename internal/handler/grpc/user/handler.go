package user

import (
	"context"

	grpc "github.com/a-h-pooladvand/microservices/api/user/v1"
	"go.uber.org/fx"
)

type GrpcHandler struct {
	grpc.UnimplementedUserServiceServer
}

type HandlerParams struct {
	fx.In
}

func NewGrpcHandler(params HandlerParams) GrpcHandler {
	return GrpcHandler{}
}

// CreateUser implements the CreateUser RPC method (legacy handler - redirect to new implementation)
func (h GrpcHandler) CreateUser(ctx context.Context, request *grpc.CreateUserRequest) (*grpc.UserResponse, error) {
	return &grpc.UserResponse{
		Id:        "",
		Name:      request.GetName(),
		Surname:   request.GetSurname(),
		CreatedAt: "",
		UpdatedAt: "",
	}, nil
}

// GetUser implements the GetUser RPC method
func (h GrpcHandler) GetUser(ctx context.Context, request *grpc.GetUserRequest) (*grpc.UserResponse, error) {
	return &grpc.UserResponse{
		Id:        request.GetId(),
		Name:      "",
		Surname:   "",
		CreatedAt: "",
		UpdatedAt: "",
	}, nil
}

// ListUsers implements the ListUsers RPC method
func (h GrpcHandler) ListUsers(ctx context.Context, request *grpc.ListUsersRequest) (*grpc.ListUsersResponse, error) {
	return &grpc.ListUsersResponse{
		Users: []*grpc.UserResponse{},
	}, nil
}

// UpdateUser implements the UpdateUser RPC method
func (h GrpcHandler) UpdateUser(ctx context.Context, request *grpc.UpdateUserRequest) (*grpc.UserResponse, error) {
	return &grpc.UserResponse{
		Id:        request.GetId(),
		Name:      request.GetName(),
		Surname:   request.GetSurname(),
		CreatedAt: "",
		UpdatedAt: "",
	}, nil
}

// DeleteUser implements the DeleteUser RPC method
func (h GrpcHandler) DeleteUser(ctx context.Context, request *grpc.DeleteUserRequest) (*grpc.DeleteUserResponse, error) {
	return &grpc.DeleteUserResponse{
		Success: true,
		Message: "user deleted",
	}, nil
}
