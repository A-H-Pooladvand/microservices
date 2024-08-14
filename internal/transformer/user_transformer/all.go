package user_transformer

import (
	"po/internal/model"
)

type AllUser struct {
	FirstName string `json:"firstName"`
}

func All(users []model.User) []AllUser {
	var result []AllUser

	for _, user := range users {
		result = append(result, AllUser{
			FirstName: user.Name,
		})
	}

	return result
}
