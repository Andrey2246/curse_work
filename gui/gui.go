package gui

import (
	"curse_work/network"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var username string

func Initialize(w fyne.Window) {
	// Username prompt
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Enter username")

	submitButton := widget.NewButton("Submit", func() {
		username = usernameEntry.Text
		if username == "" {
			dialog.ShowError(fmt.Errorf("Username cannot be empty"), w)
			return
		}
		MainMenu(w)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to File Service Client!"),
		usernameEntry,
		submitButton,
	))
}

func MainMenu(w fyne.Window) {
	uploadButton := widget.NewButton("Upload File", func() {
		UploadFileDialog(w)
	})

	downloadButton := widget.NewButton("Download File", func() {
		DownloadFileDialog(w)
	})

	wipeButton := widget.NewButton("Wipe Files", func() {
		err := network.WipeFiles()
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.ShowInformation("Success", "All files wiped from the server.", w)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Welcome, "+username),
		uploadButton,
		downloadButton,
		wipeButton,
	))
}
