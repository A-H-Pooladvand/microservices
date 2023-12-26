package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
	"net"
	"po/api/proto/ping/v1"
	"po/internal/handlers"
	"po/pkg/logger"
)

var cmd = &cobra.Command{
	Use:   "app",
	Short: "app",
	Long:  `Initializing...`,
}

func Execute() {
	logger.Boot()

	lis, err := net.Listen("tcp", ":8500")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ping.RegisterPingServiceServer(s, &handlers.Ping{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	if err = cmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
