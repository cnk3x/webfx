package db

import (
	"time"

	"github.com/cnk3x/webfx/config"
	"github.com/cnk3x/webfx/utils/log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GormOpen() (db *gorm.DB, err error) {
	dOpt := DefineConfig(Config{
		Type:     config.Get("database.type").String(),
		User:     config.Get("database.user").String(),
		Password: config.Get("database.password").String(),
		Host:     config.Get("database.host").String(),
		Port:     int(config.Get("database.port").Int()),
		Name:     config.Get("database.name").String(),
		Path:     config.Get("database.path").String(),
		Args:     config.Get("database.args").String(),
		Debug:    config.Get("database.debug").Bool(),
	})

	dialectorOpen := dialectorMap[dOpt.Type]
	if dialectorOpen == nil {
		log.Warnf("DATABASE    not support, type=%q, dns=%q", dOpt.Type, dOpt.DSN())
		return
	}

	if db, err = gorm.Open(dialectorOpen(dOpt.DSN()), DefaultConfig); err != nil {
		db.Error = err
		log.Warnf("DATABASE    can not connect, type=%q, dns=%q, err=%v", dOpt.Type, dOpt.DSN(), err)
		return
	}

	if dOpt.Debug {
		db.Logger = logger.New(LoggerWriter(true), logger.Config{SlowThreshold: 200 * time.Millisecond, LogLevel: logger.Info, Colorful: true})
		log.Debugf("DATABASE    connected, type=%q, dns=%q", dOpt.Type, dOpt.DSN())
	} else {
		db.Logger = logger.New(LoggerWriter(false), logger.Config{SlowThreshold: 200 * time.Millisecond, LogLevel: logger.Warn, Colorful: true})
		log.Infof("DATABASE    connected, type=%q", dOpt.Type)
	}

	return
}

var SingularTable = schema.NamingStrategy{SingularTable: true}

var DefaultConfig = &gorm.Config{
	NamingStrategy:                           SingularTable,
	SkipDefaultTransaction:                   true,
	PrepareStmt:                              true,
	DisableForeignKeyConstraintWhenMigrating: true,
}

type LoggerWriter bool

func (w LoggerWriter) Printf(format string, v ...any) {
	if w {
		log.Output(1, log.DEBUG, format, v...)
	} else {
		log.Output(1, log.WARN, format, v...)
	}
}
