package main

import (
	"coilforge/internal/app"
	_ "coilforge/internal/components"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
