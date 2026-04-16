package main

// File overview:
// main launches the CoilForge desktop executable.
// Subsystem: cmd entrypoint.
// It imports app for orchestration and partmanifest for side-effect catalog registration.
// Flow position: process start -> app.Run() -> app/editor/sim/render lifecycle.

import (
	"coilforge/internal/app"
	_ "coilforge/internal/partmanifest" // Registers concrete catalog parts via side effects.
	"log"
)

// main starts the application and exits on startup/runtime errors.
func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
