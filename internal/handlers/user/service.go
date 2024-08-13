package user

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s service) GetAllUsers(ctx context.Context) {
	tracer := otel.Tracer("app")
	c, span := tracer.Start(
		ctx,
		"User Servicer",
		trace.WithAttributes(attribute.String("Method", "GetAllUsers")),
	)
	defer span.End()

	s.repository.All(c)
}
