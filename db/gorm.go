package db

import (
	"gorm.io/gorm"
)

type (
	Scope       = func(tx *gorm.DB) *gorm.DB
	Transaction = func(tx *gorm.DB) error
)

// type (
// 	DeletedAt      = gorm.DeletedAt
// 	Model          = gorm.Model
// 	Association    = gorm.Association
// 	Scope          = func(tx *gorm.DB) *gorm.DB
// 	Transaction    = func(tx *gorm.DB) error
// 	OnConflict     = clause.OnConflict
// 	Expression     = clause.Expression
// 	Dialector      = gorm.Dialector
// 	Option         = gorm.Option
// 	Session        = gorm.Session
// 	Config         = gorm.Config
// 	NamingStrategy = schema.NamingStrategy
// 	Namer          = schema.Namer
// )

// var (
// 	Expr              = gorm.Expr
// 	AssignmentColumns = clause.AssignmentColumns
// 	Open              = gorm.Open

// 	loggerInfo   = logger.Info
// 	loggerSilent = logger.Silent
// 	loggerWarn   = logger.Warn

// 	_ = loggerInfo
// 	_ = loggerSilent
// 	_ = loggerWarn
// )
