package seed

import (
	"gorm.io/gorm"
	"po/internal/handlers/user"
)

type UserSeeder struct {
	model user.User
}

func (u UserSeeder) Run(db *gorm.DB) {
	var users []*user.User
	for i := 0; i < 10; i++ {
		user := user.NewUser()
		_ = user.Fake(user)
		users = append(users, user)
	}

	db.Create(users)
}
