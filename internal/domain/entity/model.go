package entity

import (
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type Model struct {
	ID        string `json:"id" faker:"-"`
	CreatedAt string `json:"created_at" faker:"-"`
	UpdatedAt string `json:"updated_at" faker:"-"`
	DeletedAt string `json:"deleted_at" faker:"-"`
}

func (m *Model) GetID() uuid.UUID {
	if m.ID == "" {
		return uuid.Nil
	}

	return uuid.MustParse(m.ID)
}

func (m *Model) Fake(v any) error {
	return faker.FakeData(v)
}
