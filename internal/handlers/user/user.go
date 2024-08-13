package user

import "po/internal/model"

type User struct {
	model.Model
	Name     string `json:"name" faker:"first_name"`
	LastName string `json:"last_name" faker:"last_name"`
}

func NewUser() *User {
	return &User{}
}
