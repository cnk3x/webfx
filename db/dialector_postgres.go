//go:build postgres

package db

import (
	"gorm.io/driver/postgres"
)

func init() {
	RegisterDialector("postgres", postgres.Open)
	RegisterDialector("postgresql", postgres.Open)
	RegisterDialector("pq", postgres.Open)
}
