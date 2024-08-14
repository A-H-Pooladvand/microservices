package user

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"po/internal/handlers/user/dto"
	"po/internal/model"
)

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s service) GetAllUsers(ctx context.Context, request dto.GetAllUsers) ([]model.User, error) {
	tracer := otel.Tracer("app")
	c, span := tracer.Start(
		ctx,
		"User Service",
		trace.WithAttributes(attribute.String("Method", "GetAllUsers")),
	)
	defer span.End()

	return s.repository.All(c, request)
}
