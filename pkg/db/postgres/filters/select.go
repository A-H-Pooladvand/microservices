package filters

import (
	"github.com/a-h-pooladvand/microservices/internal/filter"
	"gorm.io/gorm"
)

// WithSelect is a filter that can be used to select specific fields from the database.
func WithSelect(db *gorm.DB, filter filter.Filter) *gorm.DB {
	if len(filter.Select) == 0 {
		return db
	}

	return db.Select(filter.Selects())
}
