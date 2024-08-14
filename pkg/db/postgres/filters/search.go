package filters

import (
	"gorm.io/gorm"
	"po/internal/Filter"
)

// WithSearch is a filter that can be used to search the database.
func WithSearch(db *gorm.DB, filters Filter.Filter) *gorm.DB {
	if filters.Search != "" {
		return db.Where("name LIKE ?", "%"+filters.Search+"%")
	}

	return db
}
