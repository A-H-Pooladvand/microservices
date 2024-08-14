package user

import (
	"context"
	"po/internal/handlers/user/dto"
	"po/internal/model"
)

type Repository interface {
	All(ctx context.Context, request dto.GetAllUsers) ([]model.User, error)
}
