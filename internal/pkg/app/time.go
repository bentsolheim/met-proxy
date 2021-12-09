package app

import (
	"fmt"
	"github.com/palantir/stacktrace"
	"math"
	"strings"
	"time"
)

func ParseExpiresHeader(header string) (time.Time, error) {
	return ParseToEuropeOslo(header, "Mon, 2 Jan 2006 15:04:05 GMT")
}

func ParseZuluTime(zuluTime string) (time.Time, error) {
	return ParseToEuropeOslo(zuluTime, "2006-01-02T15:04:05Z")
}

func ParseToEuropeOslo(timeString string, layout string) (time.Time, error) {
	parsedTime, err := time.Parse(layout, timeString)
	if err != nil {
		return time.Time{}, stacktrace.Propagate(err, fmt.Sprintf("error while parsing string (%s)", timeString))
	}
	return ToEuropeOslo(parsedTime)
}

func ToEuropeOslo(parsedTime time.Time) (time.Time, error) {
	tz, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return time.Time{}, stacktrace.Propagate(err, "error while getting current timezone")
	}
	return parsedTime.In(tz), nil
}

func MakeReadable(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"d", days},
		{"h", hours},
		{"m", minutes},
		{"s", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		default:
			parts = append(parts, fmt.Sprintf("%d%s", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}
