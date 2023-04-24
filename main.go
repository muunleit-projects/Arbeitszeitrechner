package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	formatName     = "%-23s"
	formatZeit     = "15:04  Mon 02.01.2006"
	formatRestzeit = "%11s"
)

const (
	beginn int = iota
	standard
	max
)

type zeitraum struct {
	name  string
	dauer time.Duration
}

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

func main() {
	log.SetPrefix("Error: ")
	log.SetFlags(0)

	var anmeldung string

	if len(os.Args) > 1 {
		anmeldung = os.Args[1]
	} else {
		// if there are no command-line-arguments, assume the program to be started
		// by double-click. For the user to be able to read the output, the program
		// requests a final "enter"-command from user before closing
		defer fmt.Scanln()

		fmt.Printf("Angemeldet um [hh:mm]: ")
		_, err := fmt.Scanln(&anmeldung)
		if err != nil {
			log.Println("Benutzereingabe:", err)
			return
		}
	}

	// set time from string
	anmeldeZP, err := time.Parse("15:04", anmeldung)
	if err != nil {
		log.Println("Anmeldezeitpunkt:", err)
		return
	}

	// set date to current date
	now := time.Now()
	anmeldeZP = time.Date(
		now.Year(), now.Month(), now.Day(),
		anmeldeZP.Hour(), anmeldeZP.Minute(),
		0, 0, time.Local)

	// correct date to yesterday if workday would start after "now" otherwise
	if anmeldeZP.After(now) {
		anmeldeZP = anmeldeZP.AddDate(0, 0, -1)
	}

	// print table of times and durations until the times are reached
	var sb strings.Builder

	for i := 0; i < len(zeiten); i++ {
		fmt.Fprintf(&sb, formatName, zeiten[i].name)

		zp := anmeldeZP.Add(zeiten[i].dauer)
		sb.WriteString(zp.Format(formatZeit))
		if zp.After(now) {
			fmt.Fprintf(&sb, formatRestzeit, time.Until(zp).Round(time.Minute))
		}
		sb.WriteRune('\n')
	}
	fmt.Print(sb.String())
}
