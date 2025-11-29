package seed

import (
	"github.com/a-h-pooladvand/microservices/internal/model"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func (u UserSeeder) Run(db *gorm.DB) {
	var users []*model.User
	for i := 0; i < 10; i++ {
		user := model.NewUser()
		_ = user.Fake(user)
		users = append(users, user)
	}

	db.Create(users)
}
