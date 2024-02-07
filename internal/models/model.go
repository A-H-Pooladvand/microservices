package models

import (
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id" faker:"-"`
	CreatedAt time.Time      `json:"created_at" faker:"-"`
	UpdatedAt time.Time      `json:"updated_at" faker:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" faker:"-"`
}

func (m *Model) Fake(a any, opt ...options.OptionFunc) error {
	return faker.FakeData(a, opt...)
}
