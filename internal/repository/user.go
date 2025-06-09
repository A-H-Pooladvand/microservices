package repository

import (
	"context"
	"github.com/a-h-pooladvand/microservices/internal/domain/entity"
	"github.com/a-h-pooladvand/microservices/internal/domain/repository"
	"github.com/a-h-pooladvand/microservices/internal/dto"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var _ repository.User = (*User)(nil)

type UserParams struct {
	fx.In
	DB *gorm.DB
}

type User struct {
	db *gorm.DB
}

func NewUser(p UserParams) repository.User {
	return &User{
		db: p.DB,
	}
}

func (u User) Create(ctx context.Context, e entity.User) (*entity.User, error) {
	user := dto.ToUserModel(e)

	r := u.db.WithContext(ctx).Create(&user)

	if NotInserted(r) {
		return nil, r.Error
	}

	return dto.ToUserEntity(user), nil
}
