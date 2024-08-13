package user

import "context"

type Service interface {
	GetAllUsers(ctx context.Context)
}
