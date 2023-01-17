package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

var ColorMajorVersion = color.New(color.FgHiYellow)
var ColorSuccess = color.New(color.FgHiGreen)
var ColorInfo = color.New(color.FgHiYellow)
var ColorError = color.New(color.FgHiRed)

func DownloadWithProgress(url string, tarName string, destFolder string) (err error) {
	destTarPath := path.Join(destFolder, tarName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned status code %d", url, resp.StatusCode)
	}

	f, _ := os.OpenFile(destTarPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading",
	)
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
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

func Successf(format string, a ...interface{}) {
	format = "==> " + format
	_, _ = ColorSuccess.Printf(format, a...)
}

func Infof(format string, a ...interface{}) {
	format = "==> " + format
	_, _ = ColorInfo.Printf(format, a...)
}

func Errorf(format string, a ...interface{}) {
	format = "==> " + format
	_, _ = ColorError.Printf(format, a...)
}

func Major(a ...interface{}) {
	_, _ = ColorMajorVersion.Print(a...)
}

func Successln(a ...interface{}) {
	_, _ = ColorSuccess.Println(a...)
}

func Infoln(a ...interface{}) {
	_, _ = ColorInfo.Println(a...)
}

func Errorln(a ...interface{}) {
	_, _ = ColorError.Println(a...)
}

func Println(a ...interface{}) {
	_, _ = ColorInfo.Println(a...)
}

func Printf(format string, a ...interface{}) {
	_, _ = ColorInfo.Printf(format, a...)
}

func Print(a ...interface{}) {
	_, _ = ColorInfo.Print(a...)
}

func Fatal(a ...interface{}) {
	_, _ = ColorInfo.Print(a...)
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	_, _ = ColorInfo.Printf(format, a...)
	os.Exit(1)
}

func CheckError(err error, format string) {
	if err != nil {
		Errorf(format+": %s", err)
		os.Exit(1)
	}
}
