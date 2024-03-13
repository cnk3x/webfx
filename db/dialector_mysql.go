//go:build mysql

package db

import (
	"gorm.io/driver/mysql"
)

func init() {
	RegisterDialector("mysql", mysql.Open)
}
