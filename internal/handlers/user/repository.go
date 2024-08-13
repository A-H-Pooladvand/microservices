package user

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) All(ctx context.Context) {
	tracer := otel.Tracer("app")
	_, span := tracer.Start(
		ctx,
		"User Repository",
		trace.WithAttributes(attribute.String("Method", "All")),
	)
	defer span.End()
}
