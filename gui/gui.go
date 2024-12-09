package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MainMenu(w fyne.Window) {
	uploadButton := widget.NewButton("Encrypt File", func() {
		EncryptDialog(w)
	})

	downloadButton := widget.NewButton("Decrypt File", func() {
		DecryptDialog(w)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Welcome, User"),
		uploadButton,
		downloadButton,
	))
}
