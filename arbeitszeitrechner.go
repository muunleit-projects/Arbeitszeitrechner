package arbeitszeitrechner

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// string formats for outputting Zeitpunkt objects
const (
	nameFormat     = "%-23s"
	zeitFormat     = "15:04  Mon 02.01.2006"
	restzeitFormat = "%11s"
)

// Zeitpunkt is a struct that represents a specific moment in time.
type Zeitpunkt struct {
	time time.Time
}

// zeitraum represents a period of time with a name and duration.
type zeitraum struct {
	name  string
	dauer time.Duration
}

const (
	beginn int = iota
	standard
	max
)

var zeiten = map[int]zeitraum{
	beginn: {
		name:  "Beginn",
		dauer: 0,
	},
	standard: {
		name:  "Standard-Tag",
		dauer: time.Hour*8 + time.Minute*18,
	},
	max: {
		name:  "maximale Arbeitszeit",
		dauer: time.Hour*10 + time.Minute*45,
	},
}

// SetBeginn sets the beginning time for Zeitpunkt. It takes a string argument
// `beginn` in the format "15:04". If `beginn` is not in the correct format, it
// returns an error. `SetBeginn` then creates a new `time.Time` instance based
// on the current date and `beginn`. If `beginn` is after the current time, it
// sets the beginning time to the previous day. Finally, it sets `zp.time` to
// the computed beginning time and returns nil.
func (zp *Zeitpunkt) SetBeginn(beginn string) (err error) {
	// Parse `beginn` into a time.Time object
	startTime, err := time.Parse("15:04", beginn)
	if err != nil {
		return err
	}

	// Create a new time.Time object with the current date and `startTime` as the time
	currentTime := time.Now()
	startTime = time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		startTime.Hour(),
		startTime.Minute(),
		0,
		0,
		time.Local)

	// If the computed start time is after the current time, set it to the previous day
	if startTime.After(currentTime) {
		startTime = startTime.AddDate(0, 0, -1)
	}

	zp.time = startTime
	return nil
}

// Beginn returns the starting-time of the workday
func (zp Zeitpunkt) Beginn() string {
	return zp.time.Format(zeitFormat)
}

// Tabelle prints a table of time durations, their end times, and the time remaining until the end time.
// It writes the table to the given io.Writer.
func (zp Zeitpunkt) Tabelle(w io.Writer) {
	var table strings.Builder
	now := time.Now()

	for i := 0; i < len(zeiten); i++ {
		fmt.Fprintf(&table, nameFormat, zeiten[i].name)

		end := zp.time.Add(zeiten[i].dauer)
		table.WriteString(end.Format(zeitFormat))

		if end.After(now) {
			fmt.Fprintf(&table, restzeitFormat, time.Until(end).Round(time.Minute))
		}
		table.WriteRune('\n')
	}
	fmt.Fprint(w, table.String())
}
