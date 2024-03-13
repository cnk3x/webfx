package db

import (
	"fmt"

	"github.com/cnk3x/webfx/utils/strs"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// type (
// 	sk interface{ constraints.Integer | ~string }
// 	sv constraints.Integer
// )

func SortExpr[K strs.Int | ~string](keys []K, reverses ...bool) clause.Expression {
	reverse := len(reverses) > 0 && reverses[0]
	l := len(keys)
	expr := ``
	for i, k := range keys {
		if i > 0 {
			expr += " "
		}
		expr += fmt.Sprintf("WHEN %v THEN %d", k, lo.Ternary(reverse, l-i, i+1))
	}
	return gorm.Expr(expr)
}

func Sort[K strs.Int | ~string](db *gorm.DB, keyCol, sortCol string, keys []K, reverse ...bool) *gorm.DB {
	if wKeys := lo.Compact(lo.Uniq(keys)); len(wKeys) > 0 {
		k := clause.Column{Table: clause.CurrentTable, Name: keyCol}
		s := clause.Column{Table: clause.CurrentTable, Name: sortCol}
		v := gorm.Expr(`(CASE ? ? ELSE ? END)`, k, SortExpr(keys, reverse...), s)
		return db.Where("? IN (?)", k, wKeys).UpdateColumn(sortCol, v)
	}
	return db
}

// 满足排序字段为sort且为数值型, 单主键且为数值型
func SortSimple[K strs.Int | ~string](db *gorm.DB, keys []K, reverse ...bool) *gorm.DB {
	return Sort(db, clause.PrimaryKey, "sort", keys, reverse...)
}

func SortMapExpr[K strs.Int | ~string, S strs.Int](sorts map[K]S) clause.Expr {
	expr := ``
	for k, s := range sorts {
		if expr != "" {
			expr += " "
		}
		expr += fmt.Sprintf("WHEN %v THEN %d", k, s)
	}
	return gorm.Expr(expr)
}

func SortMap[K strs.Int | ~string, S strs.Int](db *gorm.DB, keyCol, sortCol string, sorts map[K]S) *gorm.DB {
	if keys := lo.Compact(lo.Uniq(lo.Keys(sorts))); len(keys) > 0 {
		k := clause.Column{Table: clause.CurrentTable, Name: keyCol}
		s := clause.Column{Table: clause.CurrentTable, Name: sortCol}
		v := gorm.Expr(`(CASE ? ? ELSE ? END)`, k, SortMapExpr(sorts), s)
		return db.Where("? IN (?)", k, keys).UpdateColumn(sortCol, v)
	}
	return db
}

// 满足排序字段为sort且为数值型, 单主键且为数值型
func SortMapSimple[K strs.Int | ~string, S strs.Int](db *gorm.DB, sorts map[K]S) *gorm.DB {
	return SortMap(db, clause.PrimaryKey, "sort", sorts)
}
