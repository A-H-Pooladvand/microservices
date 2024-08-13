package user

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	grpc "po/api/proto/user/v1"
	"po/pkg/trace"
)

type GrpcHandler struct {
	service Service
	tracer  trace.Tracer
	grpc.UnimplementedUserServiceServer
}

type GrpcHandlerParams struct {
	fx.In
	Service Service
	Tracer  trace.Tracer
}

func NewGrpcHandler(params GrpcHandlerParams) GrpcHandler {
	return GrpcHandler{
		service: params.Service,
		tracer:  params.Tracer,
	}
}

func (h GrpcHandler) Index(ctx context.Context, request *grpc.UserRequest) (*grpc.UserResponse, error) {
	tracer := h.tracer.FromContext(ctx)

	_, span := tracer.Start(ctx, "user grpc handler")
	defer span.End()

	span.AddEvent("Wtf")

	return &grpc.UserResponse{
		Message: fmt.Sprintf("Hello %s %s", request.GetFirstName(), request.GetLastName()),
	}, nil
}
