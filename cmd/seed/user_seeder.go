package seed

import (
	"gorm.io/gorm"
	"po/internal/model"
)

type UserSeeder struct {
	model model.User
}

func (u UserSeeder) Run(db *gorm.DB) {
	var users []*model.User
	for i := 0; i < 10; i++ {
		user := model.NewUser()
		_ = user.Fake(user)
		users = append(users, user)
	}

	db.Create(users)
}
