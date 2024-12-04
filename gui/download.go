package gui

import (
	"curse_work/encryption"
	"curse_work/network"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func DownloadFileDialog(w fyne.Window) {
	// Fetch the file list from the server
	files, err := network.ListFiles()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	if len(files) == 0 {
		dialog.ShowInformation("No Files", "There are no files available for download.", w)
		return
	}

	// Prepare options for the Select widget
	options := []string{}
	for _, file := range files {
		size := fmt.Sprintf("%d", file["size"])[12:]
		size = strings.TrimRight(size, ")")
		options = append(options, fmt.Sprintf("%s (Uploader: %s, Time: %s, Size: %s kbytes)", file["filename"], file["uploader"], file["timestamp"], size))
	}

	// Create a Select widget
	selectWidget := widget.NewSelect(options, func(selected string) {
		if selected == "" {
			return
		}

		// Extract filename
		filename := strings.Split(selected, " ")[0]

		// Download the file
		data, err := network.DownloadFile(filename)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Prompt for decryption key
		keyEntry := widget.NewEntry()
		dialog.ShowCustomConfirm("Enter Decryption Key", "OK", "Cancel", keyEntry, func(confirm bool) {
			if !confirm {
				return
			}

			key := keyEntry.Text
			if len(key) != 4 {
				dialog.ShowError(fmt.Errorf("Key must be 4 characters"), w)
				return
			}

			// Decrypt the file
			decrypted := encryption.DecryptDES(data, key)

			// Save the decrypted file
			dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil || writer == nil {
					return
				}
				defer writer.Close()

				_, err = writer.Write(decrypted)
				if err != nil {
					dialog.ShowError(err, w)
				} else {
					dialog.ShowInformation("Success", "File saved successfully.", w)
				}
			}, w)
		}, w)
	})

	// Show the Select widget in a dialog
	dialog.ShowCustom("Select File to Download", "Close", container.NewVBox(selectWidget), w)
}
