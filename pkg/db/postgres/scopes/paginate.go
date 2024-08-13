package scopes

import (
	"gorm.io/gorm"
)

// Paginate is a scope that can be used to paginate the results of a query.
func Paginate(page, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * perPage

		return db.Offset(offset).Limit(perPage)
	}
}
