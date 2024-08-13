package seed

import (
	"gorm.io/gorm"
)

type Seeder interface {
	Run(db *gorm.DB)
}
