package main

import (
	"curse_work/gui"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Create Fyne app
	a := app.NewWithID("12312312313123")
	w := a.NewWindow("ShieldCloud")

	// Launch the GUI
	gui.Initialize(w)

	w.ShowAndRun()
}
