package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const serverURL = "http://localhost:8080"

func UploadFile(filename string, data []byte, uploader string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("uploader", uploader)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	part.Write(data)
	writer.Close()

	resp, err := http.Post(serverURL+"/upload", writer.FormDataContentType(), body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to upload file" + resp.Status)
	}
	return nil
}

func DownloadFile(filename string) ([]byte, error) {
	// Prepare the URL with the filename query parameter
	resp, err := http.Get(fmt.Sprintf("%s/download?filename=%s", serverURL, filename))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the server returned an error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %s", resp.Status)
	}

	// Read the file data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ListFiles() ([]map[string]interface{}, error) {
	resp, err := http.Get(serverURL + "/list")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var files []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func WipeFiles() error {
	resp, err := http.Post(serverURL+"/wipe", "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to wipe files " + resp.Status)
	}
	return nil
}
