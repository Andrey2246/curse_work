package gui

import (
	"curse_work/encryption"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io"
)

func DecryptFile(data *[]byte, key *string, w fyne.Window) {
	dialog.ShowFileSave(func(uri fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer uri.Close()
		encrypted := encryption.EncryptDES(*data, *key)
		_, err = uri.Write(encrypted)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
	}, w) /*
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			encrypted := encryption.EncryptDES(*data, *key)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			file, err := os.Create(dir.Path() + "/" + *filename)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			defer file.Close()
			_, err = file.Write(encrypted)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Success", "File saved successfully to location: \n"+dir.Path()+*filename+".", w)
		}, w)*/
}

func DecryptDialog(w fyne.Window) {
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
		if len(key) <= 4 {
			dialog.ShowError(fmt.Errorf("Key must be at least 4 characters"), w)
			return
		}

		////////
		dialog.ShowConfirm("Plaintext location", "Enter text from keyboard?",
			func(ok bool) {
				if ok {
					//prompt for file name
					text := widget.NewMultiLineEntry()
					dialog.ShowForm("Enter file contents", "OK", "Cancel", []*widget.FormItem{{Text: "Text", Widget: text}},
						func(ok bool) {
							if ok {
								data = []byte(text.Text)
							}
							DecryptFile(&data, &key, w)
						}, w)
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
						DecryptFile(&data, &key, w)
					}, w)
				}
			}, w)
	}, w)
}