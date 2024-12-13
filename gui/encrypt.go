package gui

import (
	"crypto/rand"
	"curse_work/encryption"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io"
	"math/big"
)

func GenerateKey() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?"
	password := make([]byte, 8)
	for i := range password {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[charIndex.Int64()]
	}
	return string(password), nil
}

func EncryptFile(data *[]byte, key *string, w fyne.Window) {
	dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if uri == nil {
			dialog.ShowInformation("No File Selected", "Cancelled", w)
			return
		}
		defer uri.Close()
		encrypted := encryption.EncryptDES(*data, *key)
		for i := len(encrypted) - 1; encrypted[i] < 7; i-- {
			encrypted = encrypted[:i]
		}
		_, err = uri.Write(encrypted)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.ShowInformation("File Saved", "File saved to "+uri.URI().Path(), w)
	}, w)
}

func EncryptDialog(w fyne.Window) {
	data := make([]byte, 0)

	//prompt for key
	key := ""
	keyEntry := widget.NewEntry()
	key, err := GenerateKey()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	keyEntry.SetText(key)
	dialog.ShowCustomConfirm("Enter Encryption Key", "OK", "Cancel", keyEntry, func(confirm bool) {
		if !confirm {
			return
		}
		key := keyEntry.Text
		if len(key) <= 7 {
			dialog.ShowError(fmt.Errorf("Key must be at least 8 characters"), w)
			return
		}
		dialog.ShowConfirm("Plaintext location", "Enter text from keyboard?",
			func(ok bool) {
				if ok {
					//prompt for file name
					text := widget.NewMultiLineEntry()
					d := dialog.NewForm("Enter file contents", "OK", "Cancel", []*widget.FormItem{{Text: "Text", Widget: text}},
						func(ok bool) {
							if ok {
								data = []byte(text.Text)
							}
							EncryptFile(&data, &key, w)
						}, w)
					if data == nil {
						return
					}
					text.SetMinRowsVisible(13)
					d.Resize(fyne.NewSize(400, 400))
					d.Show()
				} else {
					dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
						if err != nil || reader == nil {
							return
						}

						data, err = io.ReadAll(reader)
						if err != nil {
							dialog.ShowError(err, w)
							return
						}
						defer reader.Close()
						if data == nil {
							return
						}
						EncryptFile(&data, &key, w)
					}, w)
				}
			}, w)
	}, w)
}
