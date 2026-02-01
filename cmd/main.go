package main

import (
	"github.com/VILJkid/current-state/ui"

	"github.com/rivo/tview"
)

func main() {
	// Create application
	app := tview.NewApplication()

	// Create menu
	menu := ui.CreateMenu(app)
	// menu.SetBorder(true)
	// menu.SetTitle("System Information")

	// Set root and run application
	if err := app.SetRoot(menu, true).Run(); err != nil {
		panic(err)
	}
}
