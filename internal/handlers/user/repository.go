package user

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"po/internal/handlers/user/dto"
	"po/internal/model"
	"po/pkg/db/postgres"
	"po/pkg/db/postgres/filters"
	"po/pkg/db/postgres/scopes"
)

type repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

func (r repository) All(ctx context.Context, request dto.GetAllUsers) ([]model.User, error) {
	tracer := otel.Tracer("app")
	_, span := tracer.Start(
		ctx,
		"User Repository",
		trace.WithAttributes(attribute.String("Method", "All")),
	)

	var users []model.User

	db := r.DB.Scopes(
		scopes.Filter(request.Filter, filters.WithSelect),
	).First(&users)

	if postgres.NotFound(db) {
		return []model.User{}, errors.New("user not found")
	}

	defer span.End()

	return users, nil
}
