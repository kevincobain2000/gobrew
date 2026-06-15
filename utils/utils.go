package utils

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"

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
		if err = body.Close(); err != nil {
			color.Errorln("==> [Error]: failed close response body", err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned status code %d", url, resp.StatusCode)
	}

	f, _ := os.OpenFile(destTarPath, os.O_CREATE|os.O_WRONLY, 0o644)
	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			color.Errorln("==> [Error]: failed close file", err.Error())
		}
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

// RemoveAll removes path and all its children, making any read-only
// directories writable first. A plain os.RemoveAll fails with
// "permission denied" when the tree contains read-only directories
// (e.g. Go module-cache style trees with mode 0555), leaving the
// directory in place. Making directories writable before removal avoids
// silently failing to delete an installed Go version.
func RemoveAll(path string) error {
	_ = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		// WalkDir does not follow symlinks and d.IsDir() is false for
		// symlinked dirs, so p stays within gobrew's version tree.
		if d.IsDir() {
			_ = os.Chmod(p, 0o755) //nolint:gosec // path is confined to the walked tree
		}
		return nil
	})
	return os.RemoveAll(path)
}

func CheckError(err error, format string) {
	if err != nil {
		color.Errorf(format+": %s", err)
		os.Exit(1)
	}
}
