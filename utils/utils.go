package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gookit/color"
	"github.com/schollz/progressbar/v3"
)

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

	defer func(body io.ReadCloser) {
		_ = body.Close()
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

func CheckError(err error, format string) {
	if err != nil {
		color.Errorf(format+": %s", err)
		os.Exit(1)
	}
}
