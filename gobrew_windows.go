package gobrew

import (
	"os"

	"github.com/kevincobain2000/gobrew/utils"
)

const (
	fileExt    = ".exe"
	tarNameExt = ".zip"
)

func removeFile(goBrewFile string) {
	goBrewOldFile := goBrewFile + ".old"
	utils.CheckError(os.Rename(goBrewFile, goBrewOldFile), "==> [Error] Cannot rename binary file")
}
