package filters

import (
	"github.com/a-h-pooladvand/microservices/internal/filter"
	"gorm.io/gorm"
)

// WithSearch is a filter that can be used to search the database.
func WithSearch(db *gorm.DB, filters filter.Filter) *gorm.DB {
	if filters.Search != "" {
		return db.Where("name LIKE ?", "%"+filters.Search+"%")
	}

	return db
}
