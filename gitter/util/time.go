package util

import (
	"fmt"
	"math"
	"time"
)

const TimeLayout = "2006-01-02T15:04:05Z"

func TimeDiff(a time.Time, b time.Time) string {
	hs := a.Sub(b).Hours()
	if hs > 24 {
		return fmt.Sprintf("%.0fd", hs/24)
	}

	hs, mf := math.Modf(hs)
	ms := mf * 60

	ms, sf := math.Modf(ms)
	ss := sf * 60

	return fmt.Sprintf("%.0fh %.0fm %.0fs", hs, ms, ss)
}
