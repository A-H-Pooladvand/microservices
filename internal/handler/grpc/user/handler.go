package user

import (
	"context"
	"fmt"
	grpc "github.com/a-h-pooladvand/microservices/api/proto/user/v1"
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

func (h GrpcHandler) Index(ctx context.Context, request *grpc.UserRequest) (*grpc.UserResponse, error) {
	//tracer := h.tracer.FromContext(ctx)
	//
	//_, span := tracer.Start(ctx, "user grpc handler")
	//defer span.End()
	//
	//span.AddEvent("Wtf")

	return &grpc.UserResponse{
		Message: fmt.Sprintf("Hello %s %s", request.GetFirstName(), request.GetLastName()),
	}, nil
}
