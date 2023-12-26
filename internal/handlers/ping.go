package handlers

import (
	"context"
	"fmt"
	"po/api/proto/ping/v1"
)

type Ping struct {
	ping.UnimplementedPingServiceServer
}

func (p Ping) Ping(_ context.Context, request *ping.PingRequest) (*ping.PingResponse, error) {
	return &ping.PingResponse{
		Message: fmt.Sprintf("Hello %s", request.Name),
	}, nil
}
