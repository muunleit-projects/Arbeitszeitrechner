package arbeitszeitrechner

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// zeitraum represents a period of time with a name and duration.
type zeitraum struct {
	name     string
	duration time.Duration
}

var zeitenraeume = []zeitraum{
	{
		name:     "Beginn",
		duration: 0,
	},
	{
		name:     "Standard-Tag",
		duration: time.Hour*8 + time.Minute*18,
	},
	{
		name:     "maximale Arbeitszeit",
		duration: time.Hour*10 + time.Minute*45,
	},
}

// Zeitpunkt is a struct that represents a specific moment in time.
type Zeitpunkt struct {
	beginn time.Time
	now    time.Time
	output io.Writer
}

type option func(*Zeitpunkt) error

// NewArbeitszeitrechner builds a new Arbeitszeitrechner-object with given
// options or setting it to standard-value if the options are omitted.
func NewArbeitszeitrechner(opts ...option) (Zeitpunkt, error) {
	z := Zeitpunkt{
		now:    time.Now(),
		output: os.Stdout,
	}

	for _, opt := range opts {
		err := opt(&z)
		if err != nil {
			return Zeitpunkt{}, err
		}
	}

	return z, nil
}

// Now sets the current time of a new Arbeitszeitrechner-object to the given
// time.
func Now(t time.Time) option {
	return func(z *Zeitpunkt) error {
		z.now = t

		return nil
	}
}

// Output sets the output of a new Arbeitszeitrechner-object to the given output.
func Output(output io.Writer) option {
	return func(z *Zeitpunkt) error {
		if output == nil {
			return errors.New("nil as io writer")
		}

		z.output = output

		return nil
	}
}

// TabelleString generates the table string without writing it to an output.
func (z *Zeitpunkt) TabelleString(checkin string) (string, error) {
	// string formats for outputting Zeitpunkt objects
	const (
		nameFormat          = "%-23s"
		timeFormat          = "15:04  Mon 02.01.2006"
		remainingTimeFormat = "%11s"
	)

	if err := z.setBeginn(checkin); err != nil {
		return "", err
	}

	var table strings.Builder

	for _, zr := range zeitenraeume {
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

	return table.String(), nil
}

// Tabelle prints a table of time durations, their end times, and the time
// remaining until the end time. It writes the table to the given io.Writer.
func (z *Zeitpunkt) Tabelle(checkin string) error {
	s, err := z.TabelleString(checkin)
	if err != nil {
		return err
	}

	fmt.Fprint(z.output, s)

	return nil
}

// Tabelle prints a table of time durations, their end times, and the time
// remaining until the end time to an io.Writer.
func Tabelle(checkin string) error {
	a, err := NewArbeitszeitrechner()
	if err != nil {
		panic("internal error")
	}

	return a.Tabelle(checkin)
}

// setBeginn sets the beginning time. It parses a string in the format "15:04"
// and sets the date to the current date. If the beginning time would be in the
// future, it reduces the date by one day.
func (z *Zeitpunkt) setBeginn(checkin string) error {
	checkinTime, err := time.Parse("15:04", checkin)
	if err != nil {
		return err
	}

	beginn := time.Date(
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
	if beginn.After(z.now) {
		beginn = beginn.AddDate(0, 0, -1)
	}

	z.beginn = beginn

	return nil
}
