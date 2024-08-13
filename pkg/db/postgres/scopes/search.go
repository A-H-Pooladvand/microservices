package scopes

import "gorm.io/gorm"

// Search is a scope that can be used to search for records based on a field.
func Search(field, query string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" LIKE ?", "%"+query+"%")
	}
}
