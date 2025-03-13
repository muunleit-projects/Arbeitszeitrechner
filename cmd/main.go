package main

import (
	"fmt"

	azr "github.com/muunleit-projects/Arbeitszeitrechner"
)

// TODO: switch to bubbletea

func main() {
	var checkin string

	fmt.Print("Eingecheckt um [hh:mm]: ")
	fmt.Scanln(&checkin)

	err := azr.Tabelle(checkin)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Bis zum nächsten Mal! (Drücken Sie Enter zum Beenden)")
	fmt.Scanln()
}
