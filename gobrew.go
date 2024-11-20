package gobrew

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/gookit/color"

	"github.com/kevincobain2000/gobrew/utils"
)

const (
	goBrewDir           string = ".gobrew"
	DefaultRegistryPath string = "https://go.dev/dl/"
	DownloadURL         string = "https://github.com/kevincobain2000/gobrew/releases/latest/download/"
	TagsAPI                    = "https://raw.githubusercontent.com/kevincobain2000/gobrew/json/golang-tags.json"
	VersionsURL         string = "https://api.github.com/repos/kevincobain2000/gobrew/releases/latest"
)

const (
	NoneVersion = "None"
	ProgramName = "gobrew"
)

// check GoBrew implement is Command interface
var _ Command = (*GoBrew)(nil)

// Command ...
type Command interface {
	ListVersions()
	ListRemoteVersions(bool) map[string][]string
	CurrentVersion() string
	Uninstall(version string)
	Install(version string) string
	Use(version string)
	Prune()
	Version(currentVersion string)
	Upgrade(currentVersion string)
	Interactive(ask bool)
}

// GoBrew struct
type GoBrew struct {
	installDir    string
	versionsDir   string
	currentDir    string
	currentBinDir string
	currentGoDir  string
	downloadsDir  string
	cacheFile     string
	Config
}

type Config struct {
	RootDir           string
	RegistryPathURL   string
	GobrewDownloadURL string
	GobrewTags        string
	GobrewVersionsURL string

	// cache settings
	TTL          time.Duration
	DisableCache bool
	ClearCache   bool
}

// NewGoBrew instance
func NewGoBrew(config Config) GoBrew {
	installDir := filepath.Join(config.RootDir, goBrewDir)
	cacheFile := filepath.Join(installDir, "cache.json")

	gb := GoBrew{
		Config:        config,
		installDir:    installDir,
		versionsDir:   filepath.Join(installDir, "versions"),
		currentDir:    filepath.Join(installDir, "current"),
		currentBinDir: filepath.Join(installDir, "current", "bin"),
		currentGoDir:  filepath.Join(installDir, "current", "go"),
		downloadsDir:  filepath.Join(installDir, "downloads"),
		cacheFile:     cacheFile,
	}

	if gb.ClearCache {
		_ = os.RemoveAll(gb.cacheFile)
	}

	return gb
}

// Interactive used by default
func (gb *GoBrew) Interactive(ask bool) {
	currentVersion := gb.CurrentVersion()
	currentMajorVersion := extractMajorVersion(currentVersion)

	latestVersion := gb.getLatestVersion()
	latestMajorVersion := extractMajorVersion(latestVersion)

	modVersion := NoneVersion
	if gb.hasModFile() {
		modVersion = gb.getModVersion()
		modVersion = extractMajorVersion(modVersion)
	}

	fmt.Println()

	if currentVersion == NoneVersion {
		color.Warnln("🚨 Installed Version", ".......", currentVersion, "⚠️")
	} else {
		var labels []string
		if modVersion != NoneVersion && currentMajorVersion != modVersion {
			labels = append(labels, "🔄 not same as go.mod")
		}
		if currentVersion != latestVersion {
			labels = append(labels, "⬆️ not latest")
		}
		label := ""
		if len(labels) > 0 {
			label = " " + color.FgRed.Render(label)
		}
		if currentVersion != latestVersion {
			color.Successln("✅ Installed Version", ".......", currentVersion+label, "\t🌟", latestVersion, "available")
		} else {
			color.Successln("✅ Installed Version", ".......", currentVersion+label, "\t🎉", "on latest")
		}
	}

	if modVersion != NoneVersion && latestMajorVersion != modVersion {
		label := " " + color.FgYellow.Render("\t⚠️  not latest")
		color.Successln("📄 go.mod Version", "   .......", modVersion+label)
	} else {
		color.Successln("📄 go.mod Version", "   .......", modVersion)
	}

	fmt.Println()

	if currentVersion == NoneVersion {
		color.Warnln("GO is not installed.")
		c := true
		if ask {
			c = askForConfirmation("Do you want to use latest GO version (" + latestVersion + ")?")
		}
		if c {
			gb.Use(latestVersion)
		}
		return
	}

	if modVersion != NoneVersion && currentMajorVersion != modVersion {
		color.Warnf("⚠️  GO Installed Version (%s) and go.mod Version (%s) are different.\n", currentMajorVersion, modVersion)
		fmt.Println("   Please consider updating your go.mod file")
		c := true
		if ask {
			c = askForConfirmation("Do you want to use GO version same as go.mod version (" + modVersion + "@latest)?")
		}
		if c {
			gb.Use(modVersion + "@latest")
		}
		return
	}

	if currentVersion != latestVersion {
		color.Warnf("GO Installed Version (%s) and GO Latest Version (%s) are different.\n", currentVersion, latestVersion)
		c := true
		if ask {
			c = askForConfirmation("Do you want to update GO to latest version (" + latestVersion + ")?")
		}
		if c {
			gb.Use(latestVersion)
		}
		return
	}
}

// Prune removes all installed versions of go except current version
func (gb *GoBrew) Prune() {
	currentVersion := gb.CurrentVersion()
	color.Infoln("==> [Info] Current version:", currentVersion)

	entries, err := os.ReadDir(gb.versionsDir)
	utils.CheckError(err, "[Error]: List versions failed")
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		utils.CheckError(err, "[Error]: List versions failed")
		files = append(files, info)
	}

	for _, f := range files {
		if f.Name() != currentVersion {
			version := f.Name()
			color.Infoln("==> [Info] Uninstalling version:", version)
			gb.Uninstall(version)
		}
	}
}

// ListVersions that are installed by dir ls
// highlight the version that is currently symbolic linked
func (gb *GoBrew) ListVersions() {
	entries, err := os.ReadDir(gb.versionsDir)
	if err != nil && os.IsNotExist(err) {
		color.Infoln("==> [Info] Nothing installed yet. Run `gobrew use latest` to install a latest version of Go.")
		return
	}

	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		utils.CheckError(err, "[Error]: List versions failed")
		files = append(files, info)
	}

	cv := gb.CurrentVersion()

	versionsSemantic := make([]*semver.Version, 0)

	for _, f := range files {
		if v, err := semver.NewVersion(f.Name()); err == nil {
			versionsSemantic = append(versionsSemantic, v)
		}
	}

	// sort semantic versions
	sort.Sort(semver.Collection(versionsSemantic))

	for _, versionSemantic := range versionsSemantic {
		version := versionSemantic.String()
		// 1.8.0 -> 1.8, if version < 1.21.0
		reMajorVersion := regexp.MustCompile("([0-9]+).([0-9]+).0")
		if len(reMajorVersion.FindStringSubmatch(version)) > 1 {
			vv, _ := strconv.Atoi(reMajorVersion.FindStringSubmatch(version)[2])
			if vv < 21 {
				if reMajorVersion.MatchString(version) {
					version = strings.Split(version, ".")[0] + "." + strings.Split(version, ".")[1]
				}
			}
		}
		if version == cv {
			version = cv + "*"
			color.Successln(version)
		} else {
			color.Infoln(version)
		}
	}

	// print rc and beta versions in the end
	for _, f := range files {
		rcVersion := f.Name()
		r := regexp.MustCompile("beta.*|rc.*")
		matches := r.FindAllString(rcVersion, -1)
		if len(matches) == 1 {
			if rcVersion == cv {
				rcVersion = cv + "*"
				color.Successln(rcVersion)
			} else {
				color.Infoln(rcVersion)
			}
		}
	}

	if cv != "" {
		color.Infoln()
		color.Infoln("current:", cv)
	}
}

// ListRemoteVersions that are installed by dir ls
func (gb *GoBrew) ListRemoteVersions(shouldPrint bool) map[string][]string {
	if shouldPrint {
		color.Infoln("==> [Info] Fetching remote versions")
	}
	tags := gb.getGolangVersions()

	var versions []string
	versions = append(versions, tags...)

	return gb.getGroupedVersion(versions, shouldPrint)
}

// CurrentVersion get current version from symb link
func (gb *GoBrew) CurrentVersion() string {
	fp, err := evalSymlinks(gb.currentBinDir)
	if err != nil {
		return NoneVersion
	}
	version := strings.TrimSuffix(fp, filepath.Join("go", "bin"))
	version = filepath.Base(version)
	if version == "." {
		return NoneVersion
	}
	return version
}

// Uninstall the given version of go
func (gb *GoBrew) Uninstall(version string) {
	if version == "" {
		color.Errorln("[Error] No version provided")
		os.Exit(1)
	}
	if gb.CurrentVersion() == version {
		color.Errorf("[Error] Version: %s you are trying to remove is your current version. Please use a different version first before uninstalling the current version\n", version)
		os.Exit(1)
	}
	gb.cleanVersionDir(version)
	color.Successf("==> [Success] Version: %s uninstalled\n", version)
}

// Install the given version of go
func (gb *GoBrew) Install(version string) string {
	if version == "" || version == NoneVersion {
		color.Errorln("[Error] No version provided")
		os.Exit(1)
	}
	version = gb.judgeVersion(version)
	if version == NoneVersion {
		color.Errorln("[Error] Version non exists")
		os.Exit(1)
	}
	if gb.existsVersion(version) {
		color.Infof("==> [Info] Version: %s exists\n", version)
		return version
	}
	gb.mkDirs(version)

	color.Infof("==> [Info] Downloading version: %s\n", version)
	gb.downloadAndExtract(version)
	gb.cleanDownloadsDir()
	color.Successf("==> [Success] Downloaded version: %s\n", version)
	return version
}

// Use a version
func (gb *GoBrew) Use(version string) {
	version = gb.Install(version)
	if gb.CurrentVersion() == version {
		color.Infof("==> [Info] Version: %s is already your current version \n", version)
		return
	}
	color.Infof("==> [Info] Changing go version to: %s \n", version)
	gb.changeSymblinkGoBin(version)
	gb.changeSymblinkGo(version)
	color.Successf("==> [Success] Changed go version to: %s\n", version)
}

// Version of GoBrew
func (gb *GoBrew) Version(currentVersion string) {
	color.Infoln("[INFO] gobrew version is", currentVersion)
}

// Upgrade of GoBrew
func (gb *GoBrew) Upgrade(currentVersion string) {
	if "v"+currentVersion == gb.getGobrewVersion() {
		color.Infoln("[INFO] your version is already newest")
		return
	}

	mkdirTemp, _ := os.MkdirTemp("", ProgramName)
	tmpFile := filepath.Join(mkdirTemp, ProgramName+fileExt)
	downloadURL, _ := url.JoinPath(gb.GobrewDownloadURL, "gobrew-"+gb.getArch()+fileExt)
	utils.CheckError(
		utils.DownloadWithProgress(downloadURL, ProgramName+fileExt, mkdirTemp),
		"[Error] Download GoBrew failed")

	source, err := os.Open(tmpFile)
	utils.CheckError(err, "[Error] Cannot open file")
	defer func(source *os.File) {
		_ = source.Close()
		utils.CheckError(os.Remove(source.Name()), "==> [Error] Cannot remove tmp file:")
	}(source)

	goBrewFile := filepath.Join(gb.installDir, "bin", ProgramName+fileExt)
	removeFile(goBrewFile)
	destination, err := os.Create(goBrewFile)
	utils.CheckError(err, "==> [Error] Cannot open file")
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)

	_, err = io.Copy(destination, source)
	utils.CheckError(err, "==> [Error] Cannot copy file")
	utils.CheckError(os.Chmod(goBrewFile, 0755), "==> [Error] Cannot set file as executable")
	color.Infoln("Upgrade successful")
}
