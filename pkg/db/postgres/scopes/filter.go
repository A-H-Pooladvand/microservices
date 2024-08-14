package scopes

import (
	"gorm.io/gorm"
	f "po/internal/Filter"
)

// filter is a function that can be used to filter the results of a query.
type filter func(db *gorm.DB, filter f.Filter) *gorm.DB

// Filter applies the filters to the database query.
func Filter(filters *f.Filter, filterable ...filter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filters == nil {
			return db
		}

		for _, callback := range filterable {
			db = callback(db, *filters)
		}

		return db
	}
}
