package utils

import (
	"io"
	"net/http"
	"os"

	"github.com/fatih/color"
)

var ColorMajorVersion = color.New(color.FgHiYellow)
var ColorSuccess = color.New(color.FgHiGreen)
var ColorInfo = color.New(color.FgHiYellow)
var ColorError = color.New(color.FgHiRed)

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

// Find takes a slice and looks for an element in it.
func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
