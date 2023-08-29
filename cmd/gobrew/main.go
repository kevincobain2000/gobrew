package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/kevincobain2000/gobrew"
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
		log.Print(usage())
		return
	}

	actionArg = args[0]
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
	}
}

func main() {
	gb := gobrew.NewGoBrew()
	switch actionArg {
	case "h", "help":
		log.Print(usage())
	case "ls", "list":
		_ = gb.ListVersions()
	case "ls-remote":
		gb.ListRemoteVersions(true)
	case "install":
		gb.Install(versionArg)
		if gb.CurrentVersion() == "" {
			gb.Use(versionArg)
		}
	case "use":
		gb.Install(versionArg)
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
`

	if runtime.GOOS == "windows" {
		msg = msg + `
    # Add gobrew to your environment variables
    PATH="%USERPROFILE%\.gobrew\current\bin;%USERPROFILE%\.gobrew\bin;%PATH%"
    GOROOT="%USERPROFILE%\.gobrew\current\go"

`
	} else {
		msg = msg + `
    # Add gobrew to your ~/.bashrc or ~/.zshrc
    export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$HOME/go/bin:$PATH"
    export GOROOT="$HOME/.gobrew/current/go"

`
	}

	return msg
}
