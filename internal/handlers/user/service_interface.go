package user

import (
	"context"
	"po/internal/handlers/user/dto"
	"po/internal/model"
)

type Service interface {
	GetAllUsers(ctx context.Context, request dto.GetAllUsers) ([]model.User, error)
}
