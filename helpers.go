package gobrew

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gookit/color"
)

func ExtractMajorVersion(version string) string {
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return ""
	}

	// Take the first two parts and join them back with a period to create the new version.
	majorVersion := strings.Join(parts[:2], ".")
	return majorVersion
}

func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s ", s)
		color.Successf("[y/n]: ")

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "" || response == "n" || response == "no" {
			return false
		}
	}
}
