package utils

import (
	"io"
	"net/http"
	"os"
)

// Download resource from url to a destination path
func Download(url string, filepath string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func BytesToString(data []byte) string {
	return string(data[:])
}
