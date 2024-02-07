package migrations

import "po/internal/models"

type Migrations []any

func Get() Migrations {
	return Migrations{
		models.User{},
	}
}
