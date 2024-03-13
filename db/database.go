package db

import (
	"fmt"
	"time"

	"github.com/cnk3x/webfx/config"
	"github.com/cnk3x/webfx/utils/log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func GormOpen() (db *gorm.DB, err error) {
	var (
		// "file::memory:"
		defDSN    = fmt.Sprintf("file::%s?_pragma=busy_timeout(3000)&_pragma=journal_mode(WAL)", config.WorkSpace("app.db"))
		dialector = config.Coalesce(config.Get("database.type").String(), "sqlite")
		dsn       = config.Coalesce(config.Get("database.dsn").String(), defDSN)
		debug     = config.Get("database.debug").Bool()
	)

	dialectorOpen := dialectorMap[dialector]
	if dialectorOpen == nil {
		log.Warnf("DATABASE    not support, type=%q, dns=%q", dialector, dsn)
		return
	}

	if db, err = gorm.Open(dialectorOpen(dsn), DefaultConfig); err != nil {
		db.Error = err
		log.Warnf("DATABASE    can not connect, type=%q, dns=%q, err=%v", dialector, dsn, err)
		return
	}

	if debug {
		db.Logger = logger.New(LoggerWriter(true), logger.Config{SlowThreshold: 200 * time.Millisecond, LogLevel: logger.Info, Colorful: true})
		log.Debugf("DATABASE    connected, type=%q, dns=%q", dialector, dsn)
	} else {
		db.Logger = logger.New(LoggerWriter(false), logger.Config{SlowThreshold: 200 * time.Millisecond, LogLevel: logger.Warn, Colorful: true})
		log.Infof("DATABASE    connected, type=%q", dialector)
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
