package durafmt

import (
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize/english"
)

const (
	Day   = 24 * time.Hour
	Week  = 7 * Day
	Month = 30 * Day
	Year  = 365 * Day
)

type section struct {
	Duration time.Duration
	Singular string
	Plural   string
}

var longSections = []section{
	{Year, "year", "years"},
	{Month, "month", "months"},
	{Week, "week", "weeks"},
	{Day, "day", "days"},
	{time.Hour, "hour", "hours"},
	{time.Minute, "minute", "minutes"},
}

// Long formats the duration in long strings.
func Long(d time.Duration) string {
	return english.OxfordWordSeries(formatStrings(d, longSections), "and")
}

var shortSections = []section{
	{Day, "d", ""},
	{time.Hour, "h", ""},
	{time.Minute, "m", ""},
}

// Short formats the duration in short strings.
func Short(d time.Duration) string {
	return strings.Join(formatStrings(d, shortSections), " ")
}

// formatStrings formats the given duration according to the given sections. It
// returns the formatted strings as well as the remainder.
func formatStrings(d time.Duration, sections []section) []string {
	var dwords = make([]string, 0, len(sections))
	var n int

	for _, section := range sections {
		n, d = divide(d, section.Duration)
		if n < 1 {
			continue
		}

		var dword string
		if section.Plural != "" {
			dword = english.Plural(n, section.Singular, section.Plural)
		} else {
			dword = fmt.Sprintf("%d%s", n, section.Singular)
		}

		dwords = append(dwords, dword)
	}

	return dwords
}

func divide(d, div time.Duration) (n int, newd time.Duration) {
	n = int(d / div)
	return n, d - time.Duration(n)*div
}
