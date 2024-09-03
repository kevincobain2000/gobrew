package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/kevincobain2000/gobrew"
	"github.com/kevincobain2000/gobrew/utils"
)

var args []string
var actionArg = ""
var versionArg = ""
var version = "dev"

var help bool
var clearCache bool
var ttl time.Duration
var disableCache bool

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
		Usage()
		return
	}

	flag.BoolVarP(&disableCache, "disable-cache", "d", false, "disable local cache")
	flag.BoolVarP(&clearCache, "clear-cache", "c", false, "clear local cache")
	flag.DurationVarP(&ttl, "ttl", "t", 20*time.Minute, "set cache duration in minutes")

	flag.BoolVarP(&help, "help", "h", false, "show usage message")

	flag.Usage = Usage

	flag.Parse()

	if help {
		Usage()
		return
	}

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
		TTL:               ttl,
		DisableCache:      disableCache,
		ClearCache:        clearCache,
	}

	gb := gobrew.NewGoBrew(config)
	switch actionArg {
	case "interactive", "info":
		gb.Interactive(true)
	case "noninteractive":
		gb.Interactive(false)
	case "h", "help":
		flag.Usage()
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

// isArgAllowed checks if the arg is allowed
// but ignored flags
//
// if the any arg is not allowed, it will return false
func isArgAllowed() bool {
	ok := true
	for i := range os.Args {
		if i == 0 {
			continue
		}
		arg := os.Args[i]
		if strings.HasPrefix(arg, "-") {
			continue
		}

		ok = Find(allowedArgs, arg)
		if !ok {
			return false
		}
	}

	return ok
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) bool {
	r := regexp.MustCompile("([0-9]+).([0-9]+)*|beta.*|rc.*|latest")
	if matches := r.FindString(val); matches != "" {
		return true
	}

	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

var Usage = func() {
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

Options:
    gobrew [--clear-cache | -c]   clear gobrew cache
    gobrew [--disable-cache | -d] disable gobrew cache
    gobrew [--ttl=20m | -t 20m]   set gobrew cache ttl, default 20m

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

	fmt.Fprintf(os.Stderr, "%s\n", msg)
}
