package grpcclient

import (
	"context"
	"fmt"
	"time"

	pb "github.com/a-h-pooladvand/microservices/api/user/v1"
	"github.com/a-h-pooladvand/microservices/internal/core/domain"
	"github.com/a-h-pooladvand/microservices/internal/core/port"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Ensure UserServiceClient implements port.UserService.
var _ port.UserService = (*UserServiceClient)(nil)

// UserServiceClient is a gRPC client for the user service.
type UserServiceClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserServiceClient creates a new UserServiceClient.
func NewUserServiceClient(address string) (*UserServiceClient, error) {
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, err
	}

	return &UserServiceClient{
		client: pb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection.
func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}

// CreateUser creates a new user via gRPC.
func (c *UserServiceClient) CreateUser(ctx context.Context, name, surname string) (*domain.User, error) {
	resp, err := c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:    name,
		Surname: surname,
	})
	if err != nil {
		return nil, mapGRPCError(err)
	}

	return protoToDomainUser(resp)
}

// GetUser retrieves a user by ID via gRPC.
func (c *UserServiceClient) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{
		Id: id.String(),
	})
	if err != nil {
		return nil, mapGRPCError(err)
	}

	return protoToDomainUser(resp)
}

// ListUsers retrieves all users with pagination via gRPC.
func (c *UserServiceClient) ListUsers(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	resp, err := c.client.ListUsers(ctx, &pb.ListUsersRequest{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		return nil, mapGRPCError(err)
	}

	users := make([]*domain.User, len(resp.Users))
	for i, u := range resp.Users {
		user, err := protoToDomainUser(u)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	return users, nil
}

// UpdateUser updates an existing user via gRPC.
func (c *UserServiceClient) UpdateUser(ctx context.Context, id uuid.UUID, name, surname string) (*domain.User, error) {
	resp, err := c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:      id.String(),
		Name:    name,
		Surname: surname,
	})
	if err != nil {
		return nil, mapGRPCError(err)
	}

	return protoToDomainUser(resp)
}

// DeleteUser deletes a user by ID via gRPC.
func (c *UserServiceClient) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := c.client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: id.String(),
	})
	if err != nil {
		return mapGRPCError(err)
	}

	return nil
}

// protoToDomainUser converts a protobuf user to a domain user.
func protoToDomainUser(resp *pb.UserResponse) (*domain.User, error) {
	id, err := uuid.Parse(resp.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, resp.GetCreatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid created_at timestamp: %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, resp.GetUpdatedAt())
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at timestamp: %w", err)
	}

	return &domain.User{
		ID:        id,
		Name:      resp.GetName(),
		Surname:   resp.GetSurname(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// mapGRPCError maps gRPC errors to domain errors.
func mapGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return domain.ErrUserNotFound
	case codes.InvalidArgument:
		return domain.ErrInvalidUser
	case codes.AlreadyExists:
		return domain.ErrUserAlreadyExists
	default:
		return domain.ErrInternal
	}
}
