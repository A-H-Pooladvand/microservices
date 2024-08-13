package routes

import (
	"google.golang.org/grpc"
	"po/api/proto/ping/v1"
	"po/api/proto/user/v1"
	"po/internal/handlers"
)

// RegisterGrpcRoutes registers gRPC routes.
func RegisterGrpcRoutes(server *grpc.Server, h *handlers.GrpcHandlers) {
	ping.RegisterPingServiceServer(server, h.Ping)
	user.RegisterUserServiceServer(server, h.User)
}
