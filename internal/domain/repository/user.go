package repository

import (
	"context"
	"github.com/a-h-pooladvand/microservices/internal/domain/entity"
)

type User interface {
	Create(ctx context.Context, e entity.User) (*entity.User, error)
}
