package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/fatih/color"
)

var ColorMajorVersion = color.New(color.FgHiYellow)
var ColorSuccess = color.New(color.FgHiGreen)
var ColorInfo = color.New(color.FgHiYellow)
var ColorError = color.New(color.FgHiRed)

func PrintDownloadPercent(done chan int64, path string, total int64) {

	var stop bool = false

	for {
		select {
		case <-done:
			stop = true
		default:

			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}

			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}

			size := fi.Size()

			if size == 0 {
				size = 1
			}

			var percent float64 = float64(size) / float64(total) * 100

			fmt.Printf("%.0f", percent)
			fmt.Println("%")
		}

		if stop {
			break
		}

		time.Sleep(time.Second)
	}
}

func Download(url string, dest string) (err error) {

	file := path.Base(url)

	ColorInfo.Printf("[Info] Downloading file %s from %s\n", file, url)

	var path bytes.Buffer
	path.WriteString(dest)
	path.WriteString("/")
	path.WriteString(file)

	start := time.Now()

	out, err := os.Create(path.String())

	if err != nil {
		return err
	}

	defer out.Close()

	headResp, err := http.Head(url)

	if err != nil {
		return err
	}

	defer headResp.Body.Close()

	size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

	if err != nil {
		return err
	}

	done := make(chan int64)

	go PrintDownloadPercent(done, path.String(), int64(size))

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	done <- n

	elapsed := time.Since(start)
	ColorInfo.Printf("[Info] Download completed in %s", elapsed)
	return err
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
