package main

import (
	"flag"
	"log"
	"os"

	"github.com/kevincobain2000/gobrew"
)

var args = []string{}
var actionArg = ""
var versionArg = ""

func init() {
	log.SetFlags(0)

	if len(os.Args) > 1 {
		if os.Args[1] == "-h" {
			log.Print(usage())
			return
		}
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
	}
}

func usage() string {
	msg := `
gobrew 1.0.2

Usage:
    gobrew help                         Show this message
    gobrew use <version>                Use <version>
    gobrew install <version>            Download and install <version> (from binary))
    gobrew uninstall <version>          Uninstall <version>
    gobrew list                         List installed versions
    gobrew ls                           Alias for list
    gobrew ls-remote                   	List remote versions

Example:
    # install and use
    gobrew use 1.16
`
	return msg
}
