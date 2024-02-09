package gobrew

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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

func symlink(oldname string, newname string) {
	utils.CheckError(
		exec.Command("cmd", "/c", "mklink", "/J", newname, oldname).Run(),
		"==> [Error]: symbolic link failed",
	)
}

// https://github.com/golang/go/issues/63703
func evalSymlinks(path string) (string, error) {
	cmd := fmt.Sprintf(
		"Get-Item -Path %s | Select-Object -ExpandProperty Target",
		path,
	)
	output, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), err
}
