package arbeitszeit

import (
	"fmt"
	"time"
)

const beginn = "Beginn:"

type zeitpunkt struct {
	name string
	zeit time.Time
}

type zeitraum struct {
	name  string
	dauer time.Duration
}

var abglZeiten []zeitraum = []zeitraum{
	{
		name:  "Standard-Tag:",
		dauer: time.Hour*8 + time.Minute*18,
	},
	{
		name:  "maximale Arbeitszeit:",
		dauer: time.Hour*10 + time.Minute*45,
	},
}

// New is a wrapper for SetBeginn
func New(s string) (zeitpunkt, error) {
	z := zeitpunkt{name: beginn}
	if err := z.SetBeginn(s); err != nil {
		return zeitpunkt{}, err
	}
	return z, nil
}

// SetBeginn returns a new zeit-object, filled up from string
// or user-input, if string is empty
func (z *zeitpunkt) SetBeginn(s string) error {
	if s == "" {
		fmt.Printf("Eingestempelt um [hh:mm]: ")
		if _, err := fmt.Scanln(&s); err != nil {
			return err
		}
		fmt.Println()
	}

	userZeit, err := time.Parse("15:04", s)
	if err != nil {
		return err
	}

	n := time.Now()

	z.zeit = time.Date(
		n.Year(), n.Month(), n.Day(),
		userZeit.Hour(), userZeit.Minute(),
		0, 0, time.Local)

	if z.zeit.After(n) {
		z.zeit = z.zeit.AddDate(0, 0, -1)
	}

	return nil
}

// Tabelle prints the times as a table
func (z zeitpunkt) Tabelle() {
	fmt.Println(z.string())

	for _, d := range abglZeiten {
		s := zeitpunkt{
			zeit: z.zeit.Add(d.dauer),
			name: d.name,
		}
		fmt.Println(s.string())
	}
}

// string is ....
func (z zeitpunkt) string() string {
	var rest string
	if z.name != beginn {
		rest = time.Until(z.zeit).Round(time.Minute).String()
	}

	s := fmt.Sprintf(
		"%-21s  %s | %8s",
		z.name,
		z.zeit.Format("15:04 | Mon 02.01.2006"),
		rest)

	return s
}
