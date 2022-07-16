package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kevincobain2000/gobrew"
)

var args = []string{}
var actionArg = ""
var versionArg = ""

var allowedArgs = []string{"h", "help", "ls", "list", "ls-remote", "install", "use", "uninstall", "self-update"}

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
		if len(versionArgSlice) == 3 && versionArgSlice[2] == "0" {
			versionArg = versionArgSlice[0] + "." + versionArgSlice[1]
		}
	}
}

func main() {
	gb := gobrew.NewGoBrew()
	switch actionArg {
	case "h", "help":
		log.Print(usage())
	case "ls", "list":
		gb.ListVersions()
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
	case "self-update":
		fmt.Println("Please execute curl cmd for self update")
		fmt.Print("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓\n\n")
		fmt.Println("curl -sLk https://git.io/gobrew | sh -")
		fmt.Print("\n↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑\n")
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
gobrew 1.6

Usage:

    gobrew use <version>           Install and use <version>
    gobrew ls                      Alias for list
    gobrew ls-remote               List remote versions (including rc|beta versions)

    gobrew install <version>       Only install <version> (binary from official or GOBREW_REGISTRY env)
    gobrew uninstall <version>     Uninstall <version>
    gobrew list                    List installed versions
    gobrew self-update             Self update this tool
    gobrew help                    Show this message

Examples:
    gobrew use 1.16         # will install and set go version to 1.16
    gobrew use 1.16.1       # will install and set go version to 1.16.1
    gobrew use 1.16rc1      # will install and set go version to 1.16rc1

    gobrew use 1.16@latest  # will install and set go version to latest version of 1.16, which is: 1.16.9
    gobrew use 1.16.x       # same as above
    gobrew use 1.16x        # same as above

    gobrew use 1.16@dev-latest   # same as @latest
                                 # will install and set go version to latest version
                                 # or beta, rc version, when major release is published

Installation Path:
    # Add gobrew to your ~/.bashrc or ~/.zshrc
    export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"

`
	return msg
}
