package seed

import (
	"po/internal/models"
	"po/pkg/postgres"
)

type UserSeeder struct {
	model models.User
}

func (u UserSeeder) Run(db *postgres.Client) {
	var users []*models.User
	for i := 0; i < 10; i++ {
		user := models.NewUser()
		_ = user.Fake(user)
		users = append(users, user)
	}

	db.Create(users)
}
