package dto

import (
	"github.com/a-h-pooladvand/microservices/internal/domain/entity"
	"github.com/a-h-pooladvand/microservices/internal/model"
	"github.com/a-h-pooladvand/microservices/internal/response"
)

func ToModelEntity(model model.Model) entity.Model {
	return entity.Model{
		ID:        model.ID.String(),
		CreatedAt: model.CreatedAt.String(),
		UpdatedAt: model.UpdatedAt.String(),
		DeletedAt: model.DeletedAt.Time.String(),
	}
}

func ToModelResponse(model entity.Model) response.Model {
	return response.Model{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: model.DeletedAt,
	}
}
