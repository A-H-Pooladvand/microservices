package dto

import (
	"github.com/a-h-pooladvand/microservices/internal/domain/entity"
	"github.com/a-h-pooladvand/microservices/internal/model"
	"github.com/a-h-pooladvand/microservices/internal/response"
)

func ToUserEntity(m *model.User) *entity.User {
	if m == nil {
		return nil
	}

	return &entity.User{
		Model:   ToModelEntity(m.Model),
		Name:    m.Name,
		Surname: m.Surname,
	}
}

func ToUserEntities(models []*model.User) []*entity.User {
	if models == nil {
		return nil
	}

	var results []*entity.User

	for _, m := range models {
		results = append(results, ToUserEntity(m))
	}

	return results
}

func ToUserModel(e entity.User) *model.User {
	return &model.User{
		Name:    e.Name,
		Surname: e.Surname,
	}
}

func ToUserResponse(m entity.User) response.User {
	return response.User{
		Model:   ToModelResponse(m.Model),
		Name:    m.Name,
		Surname: m.Surname,
	}
}
