package scopes

import "gorm.io/gorm"

// NotDeleted is a scope that can be used to query only the records that have not been soft deleted.
func NotDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}
