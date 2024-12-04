package gui

import (
	"curse_work/encryption"
	"curse_work/network"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io"
)

func UploadFileDialog(w fyne.Window) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}

		// Read the file content before closing the reader
		data, readErr := io.ReadAll(reader)
		reader.Close() // Close the file after reading
		if readErr != nil {
			dialog.ShowError(readErr, w)
			return
		}

		// Prompt for encryption key
		keyEntry := widget.NewEntry()
		dialog.ShowCustomConfirm("Enter Encryption Key", "OK", "Cancel", keyEntry, func(confirm bool) {
			if !confirm {
				return
			}

			key := keyEntry.Text
			if len(key) <= 4 {
				dialog.ShowError(fmt.Errorf("Key must be at least 4 characters"), w)
				return
			}

			// Encrypt the file
			encrypted := encryption.EncryptDES(data, key)

			// Send to server
			err := network.UploadFile(reader.URI().Name(), encrypted, username)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			dialog.ShowInformation("Success", "File uploaded successfully.", w)
		}, w)
	}, w)
}
