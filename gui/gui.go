package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Initialize(w fyne.Window) {
	/*password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter Password")
	submitButton := widget.NewButton("Submit", func() {
		sha := "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4"
		if fmt.Sprint(sha256.Sum256([]byte(password.Text))) == sha {*/
	MainMenu(w)
	/*} else {
			dialog.ShowError(errors.New("Wrong Password"), w)
		}
	})
	w.SetContent(container.NewVBox(
		widget.NewLabel("Welcome, User"),
		password,
		submitButton,
	))*/
}

func MainMenu(w fyne.Window) {
	uploadButton := widget.NewButton("Encrypt File", func() {
		EncryptDialog(w)
	})

	downloadButton := widget.NewButton("Decrypt File", func() {
		DecryptDialog(w)
	})
	visualisationButton := widget.NewButton("Visualisation of DES encryption", func() {
		Initialization(w)
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Welcome, User"),
		uploadButton,
		downloadButton,
		visualisationButton))
}
