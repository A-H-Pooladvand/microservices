package seed

import "po/pkg/postgres"

type Seeder interface {
	Run(db *postgres.Client)
}
