package model

type User struct {
	Model
	Name     string `json:"name" faker:"first_name"`
	LastName string `json:"last_name" faker:"last_name"`
}

func NewUser() *User {
	return &User{}
}
