package main

import (
	"fmt"

	"github.com/muunleit-projects/Arbeitszeitrechner/pkg/arbeitszeit"
)

func main() {
	b, err := arbeitszeit.New("")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	b.Tabelle()

	fmt.Scanln()
}
