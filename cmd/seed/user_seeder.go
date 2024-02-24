package seed

import (
	"po/internal/db"
	"po/internal/models"
)

type UserSeeder struct {
	model models.User
}

func (u UserSeeder) Run() {
	var users []*models.User
	for i := 0; i < 10; i++ {
		user := models.NewUser()
		_ = user.Fake(user)
		users = append(users, user)
	}

	db.New().Create(users)
}
