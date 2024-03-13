package db

import (
	"strings"

	"gorm.io/gorm"
)

func Paging(index, size, defSize int) Scope {
	skip, take := ParsePaging(index, size, defSize)
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Offset(skip).Limit(take)
	}
}

func Search(keywords string, columns ...string) Scope {
	return func(tx *gorm.DB) *gorm.DB {
		if keywords = strings.TrimSpace(keywords); keywords != "" && len(columns) > 0 {
			for _, col := range columns {
				tx = tx.Or(col+" like ?", "%"+keywords+"%")
			}
		}
		return tx
	}
}

func CreatedDateRange(date ...string) Scope {
	return func(tx *gorm.DB) *gorm.DB {
		var start, end int
		if len(date) == 2 {
			start, end = ParseTimeRange(date[0], date[1])
		}
		if end > 0 {
			tx = tx.Where("created_at BETWEEN ? AND ?", start, end)
		}
		return tx
	}
}
