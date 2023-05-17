package arbeitszeitrechner

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// string formats for the output
const (
	formatName     = "%-23s"
	formatZeit     = "15:04  Mon 02.01.2006"
	formatRestzeit = "%11s"
)

type Zeitpunkt struct {
	t time.Time
}

type zeitraum struct {
	name  string
	dauer time.Duration
}

const (
	beginn int = iota
	standard
	max
)

// zeiten are the names and durations of times for the output
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

// SetBeginn sets the  according to the
// beginn-parameter and the current date
func (zp *Zeitpunkt) SetBeginn(beginn string) (err error) {
	// set time from string
	zp.t, err = time.Parse("15:04", beginn)
	if err != nil {
		return err
	}
	// set date to current date
	now := time.Now()
	zp.t = time.Date(
		now.Year(), now.Month(), now.Day(),
		zp.t.Hour(), zp.t.Minute(),
		0, 0, time.Local)

	// correct date to yesterday if workday would start after "now" otherwise
	if zp.t.After(now) {
		zp.t = zp.t.AddDate(0, 0, -1)
	}
	return nil
}

// Beginn returns the starting-time of the workday
func (zp Zeitpunkt) Beginn() string {
	return zp.t.Format(formatZeit)
}

// Tabelle prints a table of times and durations until the times are reached to
// io.Writer
func (zp Zeitpunkt) Tabelle(w io.Writer) error {
	var sb strings.Builder
	now := time.Now()

	for i := 0; i < len(zeiten); i++ {
		if _, err := fmt.Fprintf(&sb, formatName, zeiten[i].name); err != nil {
			return fmt.Errorf("add to strings.Builder: %w", err)
		}

		z := zp.t.Add(zeiten[i].dauer)
		if _, err := sb.WriteString(z.Format(formatZeit)); err != nil {
			return fmt.Errorf("add to strings.Builder: %w", err)
		}
		if z.After(now) {
			if _, err := fmt.Fprintf(&sb, formatRestzeit, time.Until(z).Round(time.Minute)); err != nil {
				return fmt.Errorf("add to strings.Builder: %w", err)
			}
		}
		if _, err := sb.WriteRune('\n'); err != nil {
			return fmt.Errorf("add to strings.Builder: %w", err)
		}
	}
	fmt.Fprint(w, sb.String())
	return nil
}
