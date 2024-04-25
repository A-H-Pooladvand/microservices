package routes

import (
	"google.golang.org/grpc"
	"po/api/proto/ping/v1"
	"po/internal/handlers"
)

func RegisterGrpcRoutes(server *grpc.Server) {
	ping.RegisterPingServiceServer(server, &handlers.Ping{})
}
