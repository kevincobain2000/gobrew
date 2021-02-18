package main

import (
	"flag"
	"log"

	"github.com/kevincobain2000/gobrew"
)

var args = []string{}
var actionArg = ""
var versionArg = ""

func init() {
	log.SetFlags(0)
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
		break
	case "ls", "list":
		gb.ListVersions()
		break
	case "ls-remote":
		gb.ListRemoteVersions()
		break
	case "install":
		gb.Install(versionArg)
		if gb.CurrentVersion() == "" {
			gb.Use(versionArg)
		}
		break
	case "use":
		gb.Install(versionArg)
		gb.Use(versionArg)
		break
	case "uninstall":
		gb.Uninstall(versionArg)
		break
	}
	// gobrew.Execute()
}

func usage() string {
	msg := `
gobrew 1.0.0

Usage:
	gobrew help                         Show this message
	gobrew use <version>                Use <version>
	gobrew install <version>            Download and install <version> (from binary))
	gobrew uninstall <version>          Uninstall <version>
	gobrew list                         List installed versions
	gobrew ls                           Alias for list
	gobrew ls-remote                   	List remote versions

Example:
	# install
	gobrew install 1.16
	gobrew install 1.15.8

Reference:
	# Go versions
	https://golang.org/dl/
`
	return msg
}
