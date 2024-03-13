package db

import (
	"time"

	"github.com/cnk3x/webfx/utils/strs"

	"github.com/samber/lo"
)

func HumanDate(date time.Time) int {
	return strs.ToInt[int](date.Format("20060102150405"))
}

func ParseTimeRange(humanDate ...string) (s int, e int) {
	// 20240305//000000//14
	var start, end string
	if len(humanDate) > 0 {
		start = humanDate[0]
	}
	if len(humanDate) > 1 {
		end = humanDate[1]
	}
	start, end = strs.PadRight(start, 14, "0"), strs.PadRight(end, 14, "0")
	const timePart = 1000000
	s, e = strs.ToInt[int](start)/timePart*timePart, strs.ToInt[int](end)/timePart*timePart
	if s > e {
		s, e = e, s+timePart
	}
	return
}

func ParsePaging(index, size, defaultSize int) (skip, take int) {
	index = lo.Ternary(index < 1, 1, index)
	size = lo.Ternary(size < 1, defaultSize, size)
	skip = (index - 1) * size
	take = size
	return
}
