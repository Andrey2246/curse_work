package main

import (
	"curse_work/gui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Create Fyne app
	a := app.NewWithID("shadow2246.cursework.third.semester")
	w := a.NewWindow("ShieldCloud")

	// Launch the GUI
	gui.MainMenu(w)

	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
	w.ShowAndRun()
}
