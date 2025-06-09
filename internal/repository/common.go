package repository

import "gorm.io/gorm"

// Found checks if the transaction was successful and if any rows were affected.
func Found(tx *gorm.DB) bool {
	return tx.Error == nil && tx.RowsAffected > 0
}

// NotFound checks if the transaction was successful and if no rows were affected.
func NotFound(tx *gorm.DB) bool {
	return !Found(tx)
}

func Inserted(tx *gorm.DB) bool {
	return tx.Error == nil && tx.RowsAffected > 0
}

func NotInserted(tx *gorm.DB) bool {
	return !Inserted(tx)
}

func Affected(tx *gorm.DB) bool {
	return tx.Error == nil && tx.RowsAffected > 0
}

func NotAffected(tx *gorm.DB) bool {
	return !Affected(tx)
}
