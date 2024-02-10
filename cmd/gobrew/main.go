package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kevincobain2000/gobrew"
	"github.com/kevincobain2000/gobrew/utils"
)

var args []string
var actionArg = ""
var versionArg = ""
var version = "dev"

var allowedArgs = []string{
	"h",
	"help",
	"ls",
	"list",
	"ls-remote",
	"install",
	"use",
	"uninstall",
	"interactive",
	"noninteractive",
	"prune",
	"version",
	"self-update",
}

func init() {
	log.SetFlags(0)

	if !isArgAllowed() {
		log.Println("[Info] Invalid usage")
		log.Print(usage())
		return
	}

	flag.Parse()
	args = flag.Args()
	if len(args) == 0 {
		actionArg = "interactive"
	} else {
		actionArg = args[0]
	}

	if len(args) == 2 {
		versionArg = args[1]
		versionArgSlice := strings.Split(versionArg, ".")
		if len(versionArgSlice) == 3 {
			majorVersionNum, _ := strconv.Atoi(versionArgSlice[1])
			// Comply with: https://github.com/kevincobain2000/gobrew/issues/113
			if versionArgSlice[2] == "0" && majorVersionNum < 21 {
				// Keep complying with https://github.com/kevincobain2000/gobrew/pull/24
				versionArg = versionArgSlice[0] + "." + versionArgSlice[1]
			}
		}
		if len(versionArgSlice) == 2 {
			majorVersionNum, _ := strconv.Atoi(versionArgSlice[0])
			minorVersionNum, _ := strconv.Atoi(versionArgSlice[1])
			// Comply with: https://github.com/kevincobain2000/gobrew/issues/156
			// Check if the major version is 1 and the minor version is 21 or greater
			if majorVersionNum == 1 && minorVersionNum >= 21 {
				// Modify the versionArg to include ".0"
				versionArg += ".0"
			}
		}
	}
}

func main() {
	rootDir := os.Getenv("GOBREW_ROOT")
	if rootDir == "" {
		var err error
		rootDir, err = os.UserHomeDir()
		utils.CheckError(err, "failed get home directory and GOBREW_ROOT not defined")
	}

	registryPath := gobrew.DefaultRegistryPath
	if p := os.Getenv("GOBREW_REGISTRY"); p != "" {
		registryPath = p
	}

	config := gobrew.Config{
		RootDir:           rootDir,
		RegistryPathURL:   registryPath,
		GobrewDownloadURL: gobrew.DownloadURL,
		GobrewTags:        gobrew.TagsAPI,
		GobrewVersionsURL: gobrew.VersionsURL,
	}

	gb := gobrew.NewGoBrew(config)
	switch actionArg {
	case "interactive", "info":
		gb.Interactive(true)
	case "noninteractive":
		gb.Interactive(false)
	case "h", "help":
		log.Print(usage())
	case "ls", "list":
		gb.ListVersions()
	case "ls-remote":
		gb.ListRemoteVersions(true)
	case "install":
		gb.Install(versionArg)
		if gb.CurrentVersion() == gobrew.NoneVersion {
			gb.Use(versionArg)
		}
	case "use":
		gb.Use(versionArg)
	case "uninstall":
		gb.Uninstall(versionArg)
	case "prune":
		gb.Prune()
	case "version":
		gb.Version(version)
	case "self-update":
		gb.Upgrade(version)
	}
}

func isArgAllowed() bool {
	ok := true
	if len(os.Args) > 1 {
		_, ok = Find(allowedArgs, os.Args[1])
		if !ok {
			return false
		}
	}

	if len(os.Args) > 2 {
		_, ok = Find(allowedArgs, os.Args[1])
		if !ok {
			return false
		}
	}

	return ok
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func usage() string {
	usageMsg :=
		`
# Add gobrew to your ~/.bashrc or ~/.zshrc
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"
export GOROOT="$HOME/.gobrew/current/go"
`
	msg := `
gobrew ` + version + `

Usage:

    gobrew use <version>           Install and set <version>
    gobrew ls                      Alias for list
    gobrew ls-remote               List remote versions (including rc|beta versions)

    gobrew install <version>       Only install <version> (binary from official or GOBREW_REGISTRY env)
    gobrew uninstall <version>     Uninstall <version>
    gobrew list                    List installed versions
    gobrew self-update             Self update this tool
    gobrew prune                   Uninstall all go versions except current version
    gobrew version                 Show gobrew version
    gobrew help                    Show this message

Examples:
    gobrew use 1.16                # use go version 1.16
    gobrew use 1.16.1              # use go version 1.16.1
    gobrew use 1.16rc1             # use go version 1.16rc1

    gobrew use 1.16@latest         # use go version latest of 1.16

    gobrew use 1.16@dev-latest     # use go version latest of 1.16, including rc and beta
                                   # Note: rc and beta become no longer latest upon major release

    gobrew use mod                 # use go version listed in the go.mod file
    gobrew use latest              # use go version latest available
    gobrew use dev-latest          # use go version latest avalable, including rc and beta

Installation Path:
` + usageMsg

	return msg
}
