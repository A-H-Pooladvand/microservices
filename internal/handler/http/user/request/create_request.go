package request

import (
	"github.com/a-h-pooladvand/microservices/internal/domain/entity"
)

type Request struct {
	Name    string `json:"name" validate:"required,min=2,max=50,alphanum" example:"John"`
	Surname string `json:"surname" validate:"required,min=2,max=50,alphanum" example:"Doe"`
}

func (r Request) ToUserEntity() entity.User {
	return entity.User{
		Name:    r.Name,
		Surname: r.Surname,
	}
}
