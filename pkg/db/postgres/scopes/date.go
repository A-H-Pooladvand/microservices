package scopes

import (
	"gorm.io/gorm"
	"time"
)

// BetweenDates is a scope that can be used to query records between two dates.
func BetweenDates(field string, start, end time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" BETWEEN ? AND ?", start, end)
	}
}
