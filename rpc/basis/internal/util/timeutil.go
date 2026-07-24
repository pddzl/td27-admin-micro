package util

import "time"

// Ts formats a *time.Time as the standard DB-style timestamp string.
// Returns empty string if t is nil.
func Ts(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05.000000")
}
