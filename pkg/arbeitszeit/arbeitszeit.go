package arbeitszeit

import (
	"fmt"
	"time"
)

const (
	beginn int = iota
	standard
	max
)

const (
	formatTabelle = "%-21s  %s | %8s\n"
	formatZeit    = "15:04 | Mon 02.01.2006"
)

type zeit struct {
	t time.Time
}

type zeitraum struct {
	name  string
	dauer time.Duration
}

// zeiten are the names and durations times for the output
var zeiten map[int]zeitraum = map[int]zeitraum{
	beginn: {
		name:  "Beginn",
		dauer: 0,
	},
	standard: {
		name:  "Standard-Tag",
		dauer: time.Hour*8 + time.Minute*18,
	},
	max: {
		name:  "maximale Arbeitszeit:",
		dauer: time.Hour*10 + time.Minute*45,
	},
}

// SetBeginn sets the check-in time for the workday
func SetBeginn(s string) (zeit, error) {
	if s == "" {
		fmt.Printf("Eingestempelt um [hh:mm]: ")
		if _, err := fmt.Scanln(&s); err != nil {
			return zeit{}, err
		}
		fmt.Println()
	}

	userZeit, err := time.Parse("15:04", s)
	if err != nil {
		return zeit{}, err
	}

	n := time.Now()

	z := zeit{time.Date(
		n.Year(), n.Month(), n.Day(),
		userZeit.Hour(), userZeit.Minute(),
		0, 0, time.Local)}

	if z.t.After(n) {
		z.t = z.t.AddDate(0, 0, -1)
	}

	return z, nil
}

// Tabelle returns the list of times for the workday
func (z zeit) Tabelle() string {
	var s string
	n := time.Now()

	for i := 0; i < len(zeiten); i++ {
		zp := z.t.Add(zeiten[i].dauer)

		var until string
		if zp.After(n) {
			until = time.Until(zp).Round(time.Minute).String()
		}

		s += fmt.Sprintf(
			formatTabelle,
			zeiten[i].name,
			zp.Format(formatZeit),
			until)
	}
	return s
}
