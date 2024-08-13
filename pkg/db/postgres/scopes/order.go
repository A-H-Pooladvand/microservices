package scopes

import "gorm.io/gorm"

// OrderBy is a scope that can be used to order the results of a query.
func OrderBy(field, direction string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(field + " " + direction)
	}
}
