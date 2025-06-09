package routes

import (
	"github.com/a-h-pooladvand/microservices/api/proto/ping/v1"
	"github.com/a-h-pooladvand/microservices/api/proto/user/v1"
	"github.com/a-h-pooladvand/microservices/internal/handler"
	"google.golang.org/grpc"
)

// RegisterGrpcRoutes registers gRPC routes.
func RegisterGrpcRoutes(server *grpc.Server, h *handler.Grpc) {
	ping.RegisterPingServiceServer(server, h.Ping)
	user.RegisterUserServiceServer(server, h.User)
}
