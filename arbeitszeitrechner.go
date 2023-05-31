package arbeitszeitrechner

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// string formats for outputting Zeitpunkt objects
const (
	nameFormat          = "%-23s"
	timeFormat          = "15:04  Mon 02.01.2006"
	remainingTimeFormat = "%11s"
)

// zeitraum represents a period of time with a name and duration.
type zeitraum struct {
	name     string
	duration time.Duration
}

const (
	beginnWorkday int = iota
	standardWorkday
	maximalWorkday
)

var zeiten = map[int]zeitraum{
	beginnWorkday: {
		name:     "Beginn",
		duration: 0,
	},
	standardWorkday: {
		name:     "Standard-Tag",
		duration: time.Hour*8 + time.Minute*18,
	},
	maximalWorkday: {
		name:     "maximale Arbeitszeit",
		duration: time.Hour*10 + time.Minute*45,
	},
}

// zeitpunkt is a struct that represents a specific moment in time.
type zeitpunkt struct {
	beginn      time.Time
	currentTime time.Time
	input       io.Reader
	output      io.Writer
}

func New() *zeitpunkt {
	return &zeitpunkt{
		currentTime: time.Now(),
		output:      os.Stdout,
		input:       os.Stdin,
	}
}

func (z *zeitpunkt) SetCurrentTime(t time.Time) {
	z.currentTime = t
}

func (z *zeitpunkt) SetInput(r io.Reader) {
	z.input = r
}

func (z *zeitpunkt) SetOutput(w io.Writer) {
	z.output = w
}

// SetBeginn sets the beginning time. It takes a string from an io.Reader and
// parses it in the format "15:04" and sets the date to the current date. If the
// beginning time would be in the the future, it reduces the date by one day.
func (z *zeitpunkt) SetBeginn() error {
	// the string ist needed for time.Parse()
	var checkInString string
	if _, err := fmt.Fscanln(z.input, &checkInString); err != nil {
		return err
	}

	// Parse `checkInString` into a time.Time object
	checkInTime, err := time.Parse("15:04", checkInString)
	if err != nil {
		return err
	}

	// Create a new time.Time object with the current date and `checkInTime` as
	// the time
	z.beginn = time.Date(
		z.currentTime.Year(),
		z.currentTime.Month(),
		z.currentTime.Day(),
		checkInTime.Hour(),
		checkInTime.Minute(),
		0,
		0,
		time.Local)

	// If the computed start time is after the current time, set it to the
	// previous day
	if z.beginn.After(z.currentTime) {
		z.beginn = z.beginn.AddDate(0, 0, -1)
	}

	return nil
}

// Beginn returns the starting-time of the workday
func (zp zeitpunkt) Beginn() time.Time {
	return zp.beginn
}

// Tabelle prints a table of time durations, their end times, and the time remaining until the end time.
// It writes the table to the given io.Writer.
func (zp zeitpunkt) Tabelle(w io.Writer) {
	var table strings.Builder
	now := time.Now()

	for i := 0; i < len(zeiten); i++ {
		fmt.Fprintf(&table, nameFormat, zeiten[i].name)

		end := zp.beginn.Add(zeiten[i].duration)
		table.WriteString(end.Format(timeFormat))

		if end.After(now) {
			fmt.Fprintf(&table, remainingTimeFormat, time.Until(end).Round(time.Minute))
		}
		table.WriteRune('\n')
	}
	fmt.Fprint(w, table.String())
}
