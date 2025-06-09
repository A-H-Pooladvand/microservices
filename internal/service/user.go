package service

import (
	"context"
	"github.com/a-h-pooladvand/microservices/internal/domain/entity"
	"github.com/a-h-pooladvand/microservices/internal/domain/repository"
	"github.com/a-h-pooladvand/microservices/internal/domain/service"
	"go.uber.org/fx"
)

var _ service.User = (*User)(nil)

type UserParams struct {
	fx.In
	Repository repository.User
}

type User struct {
	repository repository.User
}

func NewUser(p UserParams) service.User {
	return &User{
		repository: p.Repository,
	}
}

func (u User) Create(ctx context.Context, e entity.User) (*entity.User, error) {
	return u.repository.Create(ctx, e)
}
