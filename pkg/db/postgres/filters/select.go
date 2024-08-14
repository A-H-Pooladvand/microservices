package filters

import (
	"gorm.io/gorm"
	"po/internal/Filter"
)

// WithSelect is a filter that can be used to select specific fields from the database.
func WithSelect(db *gorm.DB, filter Filter.Filter) *gorm.DB {
	if len(filter.Select) == 0 {
		return db
	}

	return db.Select(filter.Selects())
}
