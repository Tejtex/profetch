package utils

import (
	"fmt"
	"time"
)
func ColorText(text string, colorCode int) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", colorCode, text)
}

func Format[T any](label string, val T, colorCode int) string {
	return fmt.Sprintf("%-23s %v", ColorText(label+":", colorCode), val)
}
func FormatDuration(d time.Duration) string {
	const (
		minute = time.Minute
		hour   = time.Hour
		day    = 24 * hour
		month  = 30 * day
		year   = 365 * day
	)

	switch {
	case d >= year:
		y := d / year
		d -= y * year
		m := d / month
		return fmt.Sprintf("%dy %dm", y, m)
	case d >= month:
		m := d / month
		d -= m * month
		days := d / day
		return fmt.Sprintf("%dm %dd", m, days)
	case d >= day:
		days := d / day
		d -= days * day
		h := d / hour
		return fmt.Sprintf("%dd %dh", days, h)
	case d >= hour:
		h := d / hour
		d -= h * hour
		mins := d / minute
		return fmt.Sprintf("%dh %dmin", h, mins)
	default:
		mins := d / minute
		secs := (d - mins*minute) / time.Second
		return fmt.Sprintf("%dmin %ds", mins, secs)
	}
}