package db

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Clauses
func Clauses(conds ...clause.Expression) Scope {
	return func(d *gorm.DB) *gorm.DB { return d.Clauses(conds...) }
}

// OnConflict
var (
	DoNothing = Clauses(clause.OnConflict{DoNothing: true})
	UpdateAll = Clauses(clause.OnConflict{UpdateAll: true})
	Conflict  = func(onConflict clause.OnConflict) Scope { return Clauses(onConflict) }
	DoUpdates = func(columns ...string) Scope {
		return Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns(columns)})
	}
)
