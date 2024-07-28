package metric

import (
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
