//go:build unix

package gobrew

import (
	"os"

	"github.com/kevincobain2000/gobrew/utils"
)

const (
	fileExt    = ""
	tarNameExt = ".tar.gz"
)

func removeFile(goBrewFile string) {
	utils.CheckError(os.Remove(goBrewFile), "==> [Error] Cannot remove binary file")
}
