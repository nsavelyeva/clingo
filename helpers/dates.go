package helpers

import (
	"fmt"
	"time"
)

// GetDayMonth is a helper function to return the custom date calculated with the specified days offset.
// The result is a string in the format of "<month>-<day>" with leading zeros.
func GetDayMonth(tm time.Time, in int) string {
	if in > 0 {
		tm = tm.Add(time.Hour * 24 * time.Duration(in))
	}
	_, month, day := tm.Date()

	return fmt.Sprintf("%02d-%02d", month, day)
}
