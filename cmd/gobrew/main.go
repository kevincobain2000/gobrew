package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kevincobain2000/gobrew"
)

var args = []string{}
var actionArg = ""
var versionArg = ""

var allowedArgs = []string{"h", "help", "ls", "list", "ls-remote", "install", "use", "uninstall", "self-update"}

func init() {
	log.SetFlags(0)

	if isArgAllowed() != true {
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
		gb.ListRemoteVersions()
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
		fmt.Println("========================================")
		fmt.Println("curl -sLk https://git.io/gobrew | sh -")
		fmt.Println("========================================")
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
gobrew 1.1.0

Usage:
    gobrew help                         Show this message
    gobrew use <version>                Use <version>
    gobrew install <version>            Download and install <version> (from binary))
    gobrew uninstall <version>          Uninstall <version>
    gobrew list                         List installed versions
    gobrew ls                           Alias for list
    gobrew ls-remote                   	List remote versions
    gobrew self-update                 	Self update this tool

Example:
    # install and use
    gobrew use 1.16
`
	return msg
}
