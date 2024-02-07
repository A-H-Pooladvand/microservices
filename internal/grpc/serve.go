package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"po/api/proto/ping/v1"
	"po/internal/handlers"
)

func Serve(_ context.Context) error {
	fmt.Println("gRPC server started on [::]:8500")

	lis, err := net.Listen("tcp", ":8500")

	if err != nil {
		return err
	}

	s := grpc.NewServer()
	ping.RegisterPingServiceServer(s, &handlers.Ping{})

	err = s.Serve(lis)

	return err
}
