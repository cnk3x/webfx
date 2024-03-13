//go:build !no_sqlite

package db

import (
	"github.com/glebarez/sqlite"
)

func init() {
	RegisterDialector("sqlite", sqlite.Open)
	RegisterDialector("sqlite3", sqlite.Open)
}
