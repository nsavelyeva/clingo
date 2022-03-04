package helpers

import (
	"fmt"
	"time"
)

func GetDayMonth(tm time.Time, in int) string {
	if in > 0 {
		tm = tm.Add(time.Hour * 24 * time.Duration(in))
	}
	_, month, day := tm.Date()

	return fmt.Sprintf("%02d-%02d", month, day)
}
