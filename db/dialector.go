package db

import (
	"gorm.io/gorm"
)

var dialectorMap = map[string]func(dsn string) gorm.Dialector{}

func RegisterDialector(name string, dialector func(dsn string) gorm.Dialector) {
	dialectorMap[name] = dialector
}
