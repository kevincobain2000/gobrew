package utils

import (
	"bufio"
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
	resp, err := http.Get(url)
	if err != nil {
		ColorError.Printf("[Error]: http get file: %s \n", url)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		ColorError.Printf("[Error]: Response status code: %d \n", resp.StatusCode)
		ColorInfo.Printf("[Info]: Please wait for the file to download.\n")
	} else {
		ColorSuccess.Printf("[Success]: Response status code: %d \n", resp.StatusCode)
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath)
	wt := bufio.NewWriter(out)

	if err != nil {
		ColorError.Printf("[Error]: Creating file: %s \n", err.Error())
		return err
	}

	defer out.Close()

	_, err = io.Copy(wt, resp.Body)

	if err != nil {
		return err
	}
	wt.Flush()
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
