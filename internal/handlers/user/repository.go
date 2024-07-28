package user

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"po/pkg/postgres"
)

type Repository struct {
	DB *postgres.Client
}

func NewRepository(db *postgres.Client) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r Repository) All(ctx context.Context) {
	tracer := otel.Tracer("app")
	_, span := tracer.Start(
		ctx,
		"User Repository",
		trace.WithAttributes(attribute.String("Method", "All")),
	)
	defer span.End()
}
