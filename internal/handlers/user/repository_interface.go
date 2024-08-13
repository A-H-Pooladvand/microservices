package user

import "context"

type Repository interface {
	All(ctx context.Context)
}
