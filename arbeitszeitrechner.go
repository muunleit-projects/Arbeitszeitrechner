package arbeitszeitrechner

import (
	"errors"
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

var zeitenräume = []zeitraum{
	{name: "Beginn", duration: 0},
	{name: "Standard-Tag", duration: time.Hour*8 + time.Minute*18},
	{name: "maximale Arbeitszeit", duration: time.Hour*10 + time.Minute*45},
}

// zeitpunkt is a struct that represents a specific moment in time.
type zeitpunkt struct {
	beginn time.Time
	now    time.Time
	output io.Writer
}

type option func(*zeitpunkt) error

func NewArbeitszeitrechner(opts ...option) (zeitpunkt, error) {
	z := zeitpunkt{
		now:    time.Now(),
		output: os.Stdout,
	}

	for _, opt := range opts {
		err := opt(&z)
		if err != nil {
			return zeitpunkt{}, err
		}
	}
	return z, nil
}

func Now(t time.Time) option {
	return func(z *zeitpunkt) error {
		z.now = t
		return nil
	}
}

func Output(output io.Writer) option {
	return func(z *zeitpunkt) error {
		if output == nil {
			return errors.New("nil as io writer")
		}
		z.output = output
		return nil
	}
}

// setBeginn sets the beginning time. It parses  a string in the format "15:04"
// and sets the date to the current date. If the beginning time would be in the
// the future, it reduces the date by one day.
func (z *zeitpunkt) setBeginn(checkin string) error {
	checkinTime, err := time.Parse("15:04", checkin)
	if err != nil {
		return err
	}

	z.beginn = time.Date(
		z.now.Year(),
		z.now.Month(),
		z.now.Day(),
		checkinTime.Hour(),
		checkinTime.Minute(),
		0,
		0,
		time.Local)

	// If the computed start time is after the current time, set it to the
	// previous day
	if z.beginn.After(z.now) {
		z.beginn = z.beginn.AddDate(0, 0, -1)
	}

	return nil
}

// Tabelle prints a table of time durations, their end times, and the time
// remaining until the end time. It writes the table to the given io.Writer.
func (z *zeitpunkt) Tabelle(checkin string) error {
	if err := z.setBeginn(checkin); err != nil {
		return err
	}

	var table strings.Builder

	// for i := 0; i < len(zeiten); i++ {
	for _, zr := range zeitenräume {
		fmt.Fprintf(&table, nameFormat, zr.name)

		end := z.beginn.Add(zr.duration)
		table.WriteString(end.Format(timeFormat))

		if end.After(z.now) {
			fmt.Fprintf(&table,
				remainingTimeFormat,
				end.Sub(z.now).Round(time.Minute))
		}
		table.WriteRune('\n')
	}
	fmt.Fprint(z.output, table.String())
	return nil
}

// Tabelle prints a table of time durations, their end times, and the time
// remaining until the end time. It writes the table to the given io.Writer.
func Tabelle(checkin string) error {
	a, err := NewArbeitszeitrechner()
	if err != nil {
		panic("internal error")
	}
	return a.Tabelle(checkin)
}
