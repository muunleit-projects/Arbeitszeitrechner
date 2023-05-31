package main

import (
	//"fmt"
	"fmt"
	"log"
	//"os"
	azr "github.com/muunleit-projects/Arbeitszeitrechner"
)

func main() {
	log.SetPrefix("Error: ")
	log.SetFlags(0)

	arbeit := azr.New()
	err := arbeit.SetBeginn()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(arbeit.Beginn())

	/*
		var anmeldung string

		if len(os.Args) > 1 {
			anmeldung = os.Args[1]
		} else {
			// if there are no command-line-arguments, assume the program to be started
			// by double-click. For the user to be able to read the output, the program
			// requests a final "enter"-command from the user before closing
			defer fmt.Scanln()

			fmt.Printf("Angemeldet um [hh:mm]: ")
			if _, err := fmt.Scanln(&anmeldung); err != nil {
				log.Println("Benutzereingabe:", err)
				return
			}
		}

		var zeitpunkt az.Zeitpunkt

		if err := zeitpunkt.SetBeginn(anmeldung); err != nil {
			log.Println(err)
			return
		}

		if err := zeitpunkt.Tabelle(os.Stdout); err != nil {
			log.Println(err)
			return
		}
	*/
}
