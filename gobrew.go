package gobrew

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/kevincobain2000/gobrew/utils"
)

const (
	goBrewDir     string = ".gobrew"
	registryPath  string = "https://golang.org/dl/"
	fetchTagsRepo string = "https://github.com/golang/go"
)

// Command ...
type Command interface {
	ListVersions()
	ListRemoteVersions()
	CurrentVersion() string
	Uninstall(version string)
	Install(version string)
	Use(version string)
	Helper
}

// GoBrew struct
type GoBrew struct {
	homeDir       string
	installDir    string
	versionsDir   string
	currentDir    string
	currentBinDir string
	currentGoDir  string
	downloadsDir  string
	Command
}

// Helper ...
type Helper interface {
	getArch() string
	existsVersion(version string) bool
	cleanVersionDir(version string)
	mkdirs(version string)
	getVersionDir(version string) string
	downloadAndExtract(version string)
	changeSymblinkGoBin(version string)
	changeSymblinkGo(version string)
}

var gb GoBrew

// NewGoBrew instance
func NewGoBrew() GoBrew {
	gb.homeDir = os.Getenv("HOME")
	gb.installDir = filepath.Join(gb.homeDir, goBrewDir)
	gb.versionsDir = filepath.Join(gb.installDir, "versions")
	gb.currentDir = filepath.Join(gb.installDir, "current")
	gb.currentBinDir = filepath.Join(gb.installDir, "current", "bin")
	gb.currentGoDir = filepath.Join(gb.installDir, "current", "go")
	gb.downloadsDir = filepath.Join(gb.installDir, "downloads")

	return gb
}

func (gb *GoBrew) getArch() string {
	return runtime.GOOS + "-" + runtime.GOARCH
}

// ListVersions that are installed by dir ls
// highlight the version that is currently symbolic linked
func (gb *GoBrew) ListVersions() {
	files, err := ioutil.ReadDir(gb.versionsDir)
	if err != nil {
		utils.ColorError.Printf("[Error]: List versions failed: %s", err)
		os.Exit(0)
	}
	cv := gb.CurrentVersion()

	versionsSemantic := make([]*semver.Version, 0)

	for _, f := range files {
		v, err := semver.NewVersion(f.Name())
		if err != nil {
			// utils.ColorError.Printf("Error parsing version: %s", err)
		} else {
			versionsSemantic = append(versionsSemantic, v)
		}
	}

	// sort semantic versions
	sort.Sort(semver.Collection(versionsSemantic))

	for _, versionSemantic := range versionsSemantic {
		version := versionSemantic.String()
		// 1.8.0 -> 1.8
		reMajorVersion, _ := regexp.Compile("[0-9]+.[0-9]+.0")
		if reMajorVersion.MatchString((version)) {
			version = strings.Split(version, ".")[0] + "." + strings.Split(version, ".")[1]
		}

		if version == cv {
			version = cv + "*"
			utils.ColorSuccess.Println(version)
		} else {
			log.Println(version)
		}
	}

	// print rc and beta versions in the end
	for _, f := range files {
		rcVersion := f.Name()
		r, _ := regexp.Compile("beta.*|rc.*")
		matches := r.FindAllString(rcVersion, -1)
		if len(matches) == 1 {
			if rcVersion == cv {
				rcVersion = cv + "*"
				utils.ColorSuccess.Println(rcVersion)
			} else {
				log.Println(rcVersion)
			}
		}
	}

	if cv != "" {
		log.Println()
		log.Printf("current: %s", cv)
	}
}

// ListRemoteVersions that are installed by dir ls
func (gb *GoBrew) ListRemoteVersions() {
	log.Println("[Info]: Fetching remote versions")
	cmd := exec.Command(
		"git",
		"ls-remote",
		// "--sort=version:refname",
		"--tags",
		fetchTagsRepo,
		"go*")
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.ColorError.Printf("[Error]: List remote versions failed: %s", err)
		os.Exit(0)
	}
	tagsRaw := utils.BytesToString(output)
	r, _ := regexp.Compile("tags/go.*")

	matches := r.FindAllString(tagsRaw, -1)
	versions := make([]string, len(matches))
	for _, match := range matches {
		versionTag := strings.ReplaceAll(match, "tags/go", "")
		versions = append(versions, versionTag)
	}
	printGroupedVersions(versions)
}

func printGroupedVersions(versions []string) {
	groupedVersions := make(map[string][]string)
	for _, version := range versions {
		parts := strings.Split(version, ".")
		if len(parts) > 1 {
			majorVersion := fmt.Sprintf("%s.%s", parts[0], parts[1])
			r, _ := regexp.Compile("beta.*|rc.*")
			matches := r.FindAllString(majorVersion, -1)
			if len(matches) == 1 {
				majorVersion = strings.Split(version, matches[0])[0]
			}
			groupedVersions[majorVersion] = append(groupedVersions[majorVersion], version)
		}
	}

	// groupedVersionKeys := []string{"1", "1.1", "1.2", ..., "1.17"}
	groupedVersionKeys := make([]string, 0, len(groupedVersions))
	for groupedVersionKey := range groupedVersions {
		groupedVersionKeys = append(groupedVersionKeys, groupedVersionKey)
	}

	versionsSemantic := make([]*semver.Version, 0)
	for _, r := range groupedVersionKeys {
		v, err := semver.NewVersion(r)
		if err != nil {
			// utils.ColorError.Printf("Error parsing version: %s", err)
		} else {
			versionsSemantic = append(versionsSemantic, v)
		}
	}

	// sort semantic versions
	sort.Sort(semver.Collection(versionsSemantic))

	// match 1.0.0 or 2.0.0
	reTopVersion, _ := regexp.Compile("[0-9]+.0.0")

	for _, versionSemantic := range versionsSemantic {
		strKey := versionSemantic.String()
		lookupKey := ""
		versionParts := strings.Split(strKey, ".")

		// prepare lookup key for the grouped version map.
		// 1.0.0 -> 1.0, 1.1.1 -> 1.1
		lookupKey = versionParts[0] + "." + versionParts[1]
		// On match 1.0.0, print 1. On match 2.0.0 print 2
		if reTopVersion.MatchString((strKey)) {
			utils.ColorMajorVersion.Print(versionParts[0])
			fmt.Print("\t")
		} else {
			utils.ColorMajorVersion.Print(lookupKey)
			fmt.Print("\t")
		}

		groupedVersionsSemantic := make([]*semver.Version, 0)
		for _, r := range groupedVersions[lookupKey] {
			v, err := semver.NewVersion(r)
			if err != nil {
				// utils.ColorError.Printf("Error parsing version: %s", err)
			} else {
				groupedVersionsSemantic = append(groupedVersionsSemantic, v)
			}

		}
		// sort semantic versions
		sort.Sort(semver.Collection(groupedVersionsSemantic))

		for _, gvSemantic := range groupedVersionsSemantic {
			fmt.Print(gvSemantic.String() + "  ")
		}

		// print rc and beta versions in the end
		for _, rcVersion := range groupedVersions[lookupKey] {
			r, _ := regexp.Compile("beta.*|rc.*")
			matches := r.FindAllString(rcVersion, -1)
			if len(matches) == 1 {
				fmt.Print(rcVersion + "  ")
			}
		}
		fmt.Println()
	}
}

func (gb *GoBrew) existsVersion(version string) bool {
	path := filepath.Join(gb.versionsDir, version, "go")
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// CurrentVersion get current version from symb link
func (gb *GoBrew) CurrentVersion() string {

	fp, err := filepath.EvalSymlinks(gb.currentBinDir)
	if err != nil {
		return ""
	}

	version := strings.ReplaceAll(fp, "/go/bin", "")
	version = strings.ReplaceAll(version, gb.versionsDir, "")
	version = strings.ReplaceAll(version, "/", "")
	return version
}

// Uninstall the given version of go
func (gb *GoBrew) Uninstall(version string) {
	if version == "" {
		log.Fatal("[Error] No version provided")
	}
	if gb.CurrentVersion() == version {
		utils.ColorError.Printf("[Error] Version: %s you are trying to remove is your current version. Please use a different version first before uninstalling the current version\n", version)
		os.Exit(0)
		return
	}
	if !gb.existsVersion(version) {
		utils.ColorError.Printf("[Error] Version: %s you are trying to remove is not installed\n", version)
		os.Exit(0)
	}
	gb.cleanVersionDir(version)
	utils.ColorSuccess.Printf("[Success] Version: %s uninstalled\n", version)
}

func (gb *GoBrew) cleanVersionDir(version string) {
	os.RemoveAll(gb.getVersionDir(version))
}

func (gb *GoBrew) cleanDownloadsDir() {
	os.RemoveAll(gb.downloadsDir)
}

// Install the given version of go
func (gb *GoBrew) Install(version string) {
	if version == "" {
		log.Fatal("[Error] No version provided")
	}
	gb.mkdirs(version)
	if gb.existsVersion(version) {
		utils.ColorInfo.Printf("[Info] Version: %s exists \n", version)
		return
	}

	utils.ColorInfo.Printf("[Info] Downloading version: %s \n", version)
	gb.downloadAndExtract(version)
	gb.cleanDownloadsDir()
	utils.ColorSuccess.Printf("[Success] Downloaded version: %s\n", version)
}

// Use a version
func (gb *GoBrew) Use(version string) {
	if gb.CurrentVersion() == version {
		utils.ColorInfo.Printf("[Info] Version: %s is already your current version \n", version)
		return
	}
	utils.ColorInfo.Printf("[Info] Changing go version to: %s \n", version)
	gb.changeSymblinkGoBin(version)
	gb.changeSymblinkGo(version)
	utils.ColorSuccess.Printf("[Success] Changed go version to: %s\n", version)
}

func (gb *GoBrew) mkdirs(version string) {
	os.MkdirAll(gb.installDir, os.ModePerm)
	os.MkdirAll(gb.currentDir, os.ModePerm)
	os.MkdirAll(gb.versionsDir, os.ModePerm)
	os.MkdirAll(gb.getVersionDir(version), os.ModePerm)
	os.MkdirAll(gb.downloadsDir, os.ModePerm)
}

func (gb *GoBrew) getVersionDir(version string) string {
	return filepath.Join(gb.versionsDir, version)
}
func (gb *GoBrew) downloadAndExtract(version string) {
	tarName := "go" + version + "." + gb.getArch() + ".tar.gz"

	downloadURL := registryPath + tarName
	utils.ColorInfo.Printf("[Info] Downloading from: %s \n", downloadURL)

	err := utils.Download(
		downloadURL,
		filepath.Join(gb.downloadsDir, tarName))

	if err != nil {
		gb.cleanVersionDir(version)
		utils.ColorInfo.Printf("[Info]: Downloading version failed: %s \n", err)
		utils.ColorError.Printf("[Error]: Please check connectivity to url: %s\n", downloadURL)
		os.Exit(0)
	}

	cmd := exec.Command(
		"tar",
		"-xf",
		filepath.Join(gb.downloadsDir, tarName),
		"-C",
		gb.getVersionDir(version))

	utils.ColorInfo.Printf("[Success] Untar to %s\n", gb.getVersionDir(version))
	_, err = cmd.Output()
	if err != nil {
		// clean up dir
		gb.cleanVersionDir(version)
		utils.ColorInfo.Printf("[Info]: Untar failed: %s \n", err)
		utils.ColorError.Printf("[Error]: Please check if version exists from url: %s\n", downloadURL)
		os.Exit(0)
	}
}

func (gb *GoBrew) changeSymblinkGoBin(version string) {

	goBinDst := filepath.Join(gb.versionsDir, version, "/go/bin")
	os.RemoveAll(gb.currentBinDir)

	cmd := exec.Command("ln", "-snf", goBinDst, gb.currentBinDir)

	_, err := cmd.Output()
	if err != nil {
		utils.ColorError.Printf("[Error]: symbolic link failed: %s\n", err)
		os.Exit(0)
	}

}
func (gb *GoBrew) changeSymblinkGo(version string) {

	os.RemoveAll(gb.currentGoDir)
	versionGoDir := filepath.Join(gb.versionsDir, gb.CurrentVersion(), "go")
	cmd := exec.Command("ln", "-snf", versionGoDir, gb.currentGoDir)

	_, err := cmd.Output()
	if err != nil {
		utils.ColorError.Printf("[Error]: symbolic link failed: %s\n", err)
		os.Exit(0)
	}
}
