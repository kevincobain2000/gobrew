//go:build unix

package gobrew

import (
	"os"
	"path/filepath"

	"github.com/kevincobain2000/gobrew/utils"
)

const (
	fileExt    = ""
	tarNameExt = ".tar.gz"
)

func removeFile(goBrewFile string) {
	utils.CheckError(os.Remove(goBrewFile), "==> [Error] Cannot remove binary file")
}

func symlink(oldname string, newname string) {
	utils.CheckError(os.Symlink(oldname, newname), "==> [Error]: symbolic link failed")
}

func evalSymlinks(path string) (string, error) {
	fp, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	return fp, nil
}
