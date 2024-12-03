package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create a new Fyne application
	myApp := app.New()
	myWindow := myApp.NewWindow("File Uploader")

	// UI elements
	statusLabel := widget.NewLabel("Select a file to upload")
	fileEntry := widget.NewEntry()
	fileEntry.SetPlaceHolder("Choose file path...")
	chooseButton := widget.NewButton("Choose File", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				statusLabel.SetText("Error opening file")
				return
			}
			if reader != nil {
				fileEntry.SetText(reader.URI().Path())
			}
		}, myWindow)
	})
	uploadButton := widget.NewButton("Upload", func() {
		filePath := fileEntry.Text
		if filePath == "" {
			statusLabel.SetText("Please select a file first.")
			return
		}

		// Upload file
		err := uploadFile("http://your-server-url/upload", filePath)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Upload failed: %v", err))
		} else {
			statusLabel.SetText("File uploaded successfully!")
		}
	})

	// Layout
	content := container.NewVBox(
		statusLabel,
		fileEntry,
		chooseButton,
		uploadButton,
	)
	myWindow.SetContent(content)
	// Run the application
	myWindow.ShowAndRun()
}

// uploadFile uploads a file to the server
func uploadFile(serverURL, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Create a multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form file field
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return fmt.Errorf("could not create form file: %w", err)
	}

	// Copy the file content into the form field
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}

	// Close the multipart writer
	writer.Close()

	// Create a POST request
	req, err := http.NewRequest("POST", serverURL, body)
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with: %s", resp.Status)
	}

	return nil
}
