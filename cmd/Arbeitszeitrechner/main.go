package main

import (
	"fmt"
	"os"

	"github.com/muunleit-projects/Arbeitszeitrechner/pkg/arbeitszeit"
)

func main() {
	var arg string

	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	b, err := arbeitszeit.SetBeginn(arg)
	if err != nil {
		fmt.Println("Error:", err)
		return

	}

	b.Tabelle()

	fmt.Scanln()
}
