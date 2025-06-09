package module

import (
	grpc "github.com/a-h-pooladvand/microservices/internal/handler/grpc/user"
	http "github.com/a-h-pooladvand/microservices/internal/handler/http/user"
	"github.com/a-h-pooladvand/microservices/internal/repository"
	"github.com/a-h-pooladvand/microservices/internal/service"
	"go.uber.org/fx"
)

var User = fx.Module("user", fx.Provide(
	http.NewHttpHandler,
	grpc.NewGrpcHandler,
	repository.NewUser,
	service.NewUser,
))
