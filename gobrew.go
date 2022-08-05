package gobrew

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/c4milo/unpackit"
	"github.com/kevincobain2000/gobrew/utils"
)

const (
	goBrewDir           string = ".gobrew"
	defaultRegistryPath string = "https://golang.org/dl/"
	fetchTagsRepo       string = "https://github.com/golang/go"
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
func (gb *GoBrew) ListVersions() error {
	entries, err := os.ReadDir(gb.versionsDir)
	if err != nil {
		_, _ = utils.ColorError.Printf("[Error]: List versions failed: %s", err)
		os.Exit(1)
	}
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			_, _ = utils.ColorError.Printf("[Error]: List versions failed: %s", err)
			os.Exit(1)
		}
		files = append(files, info)
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
		if reMajorVersion.MatchString(version) {
			version = strings.Split(version, ".")[0] + "." + strings.Split(version, ".")[1]
		}

		if version == cv {
			version = cv + "*"
			_, _ = utils.ColorSuccess.Println(version)
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
				_, _ = utils.ColorSuccess.Println(rcVersion)
			} else {
				log.Println(rcVersion)
			}
		}
	}

	if cv != "" {
		log.Println()
		log.Printf("current: %s", cv)
	}
	return nil
}

// ListRemoteVersions that are installed by dir ls
func (gb *GoBrew) ListRemoteVersions(print bool) map[string][]string {
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
		_, _ = utils.ColorError.Printf("[Error]: List remote versions failed: %s", err)
		os.Exit(1)
	}
	tagsRaw := utils.BytesToString(output)
	r, _ := regexp.Compile("tags/go.*")

	matches := r.FindAllString(tagsRaw, -1)
	versions := make([]string, len(matches))
	for _, match := range matches {
		versionTag := strings.ReplaceAll(match, "tags/go", "")
		versions = append(versions, versionTag)
	}
	return gb.getGroupedVersion(versions, print)
}

func (gb *GoBrew) getGroupedVersion(versions []string, print bool) map[string][]string {
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
		maxPerLine := 0
		strKey := versionSemantic.String()
		lookupKey := ""
		versionParts := strings.Split(strKey, ".")

		// prepare lookup key for the grouped version map.
		// 1.0.0 -> 1.0, 1.1.1 -> 1.1
		lookupKey = versionParts[0] + "." + versionParts[1]
		// On match 1.0.0, print 1. On match 2.0.0 print 2
		if reTopVersion.MatchString(strKey) {
			if print {
				_, _ = utils.ColorMajorVersion.Print(versionParts[0])
			}
			gb.print("\t", print)
		} else {
			if print {
				_, _ = utils.ColorMajorVersion.Print(lookupKey)
			}
			gb.print("\t", print)
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
			maxPerLine++
			if maxPerLine == 6 {
				maxPerLine = 0
				gb.print("\n\t", print)
			}
			gb.print(gvSemantic.String()+"  ", print)
		}

		// print rc and beta versions in the end
		for _, rcVersion := range groupedVersions[lookupKey] {
			r, _ := regexp.Compile("beta.*|rc.*")
			matches := r.FindAllString(rcVersion, -1)
			if len(matches) == 1 {
				gb.print(rcVersion+"  ", print)
				maxPerLine++
				if maxPerLine == 6 {
					maxPerLine = 0
					gb.print("\n\t", print)
				}
			}
		}
		gb.print("\n", print)
		gb.print("\n", print)
	}
	return groupedVersions
}

func (gb *GoBrew) print(message string, shouldPrint bool) {
	if shouldPrint {
		fmt.Print(message)
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
		_, _ = utils.ColorError.Printf("[Error] Version: %s you are trying to remove is your current version. Please use a different version first before uninstalling the current version\n", version)
		os.Exit(1)
		return
	}
	if !gb.existsVersion(version) {
		_, _ = utils.ColorError.Printf("[Error] Version: %s you are trying to remove is not installed\n", version)
		os.Exit(1)
	}
	gb.cleanVersionDir(version)
	_, _ = utils.ColorSuccess.Printf("[Success] Version: %s uninstalled\n", version)
}

func (gb *GoBrew) cleanVersionDir(version string) {
	_ = os.RemoveAll(gb.getVersionDir(version))
}

func (gb *GoBrew) cleanDownloadsDir() {
	_ = os.RemoveAll(gb.downloadsDir)
}

// Install the given version of go
func (gb *GoBrew) Install(version string) {
	if version == "" {
		log.Fatal("[Error] No version provided")
	}
	version = gb.judgeVersion(version)
	gb.mkdirs(version)
	if gb.existsVersion(version) {
		_, _ = utils.ColorInfo.Printf("[Info] Version: %s exists \n", version)
		return
	}

	_, _ = utils.ColorInfo.Printf("[Info] Downloading version: %s \n", version)
	gb.downloadAndExtract(version)
	gb.cleanDownloadsDir()
	_, _ = utils.ColorSuccess.Printf("[Success] Downloaded version: %s\n", version)
}

func (gb *GoBrew) judgeVersion(version string) string {
	judgedVersion := ""
	rcBetaOk := false
	reRcOrBeta, _ := regexp.Compile("beta.*|rc.*")
	// check if version string ends with x

	if strings.HasSuffix(version, "x") {
		judgedVersion = version[:len(version)-1]
	}

	if strings.HasSuffix(version, ".x") {
		judgedVersion = version[:len(version)-2]
	}
	if strings.HasSuffix(version, "@latest") {
		judgedVersion = version[:len(version)-7]
	}
	if strings.HasSuffix(version, "@dev-latest") {
		judgedVersion = version[:len(version)-11]
		rcBetaOk = true
	}

	if version == "latest" || version == "dev-latest" {
		groupedVersions := gb.ListRemoteVersions(false) // donot print
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
		// loop in reverse
		for i := len(versionsSemantic) - 1; i >= 0; i-- {
			judgedVersions := groupedVersions[versionsSemantic[i].Original()]
			// get last element
			if version == "dev-latest" {
				return judgedVersions[len(judgedVersions)-1]
			}

			// loop in reverse
			for j := len(judgedVersions) - 1; j >= 0; j-- {
				matches := reRcOrBeta.FindAllString(judgedVersions[j], -1)
				if len(matches) == 0 {
					return judgedVersions[j]
				}
			}
		}

		latest := versionsSemantic[len(versionsSemantic)-1].String()
		return gb.judgeVersion(latest)
	}

	if judgedVersion != "" {
		groupedVersions := gb.ListRemoteVersions(false) // donot print
		// check if judgedVersion is in the groupedVersions
		if _, ok := groupedVersions[judgedVersion]; ok {
			// get last item in the groupedVersions excluding rc and beta
			// loop in reverse groupedVersions
			for i := len(groupedVersions[judgedVersion]) - 1; i >= 0; i-- {
				matches := reRcOrBeta.FindAllString(groupedVersions[judgedVersion][i], -1)
				if len(matches) == 0 {
					return groupedVersions[judgedVersion][i]
				}
			}
			if rcBetaOk {
				// return last element including beta and rc if present
				return groupedVersions[judgedVersion][len(groupedVersions[judgedVersion])-1]
			}
		}
	}

	return version
}

// Use a version
func (gb *GoBrew) Use(version string) {
	version = gb.judgeVersion(version)
	if gb.CurrentVersion() == version {
		_, _ = utils.ColorInfo.Printf("[Info] Version: %s is already your current version \n", version)
		return
	}
	_, _ = utils.ColorInfo.Printf("[Info] Changing go version to: %s \n", version)
	gb.changeSymblinkGoBin(version)
	gb.changeSymblinkGo(version)
	_, _ = utils.ColorSuccess.Printf("[Success] Changed go version to: %s\n", version)
}

func (gb *GoBrew) mkdirs(version string) {
	_ = os.MkdirAll(gb.installDir, os.ModePerm)
	_ = os.MkdirAll(gb.currentDir, os.ModePerm)
	_ = os.MkdirAll(gb.versionsDir, os.ModePerm)
	_ = os.MkdirAll(gb.getVersionDir(version), os.ModePerm)
	_ = os.MkdirAll(gb.downloadsDir, os.ModePerm)
}

func (gb *GoBrew) getVersionDir(version string) string {
	return filepath.Join(gb.versionsDir, version)
}
func (gb *GoBrew) downloadAndExtract(version string) {
	tarName := "go" + version + "." + gb.getArch() + ".tar.gz"

	registryPath := defaultRegistryPath
	if p := os.Getenv("GOBREW_REGISTRY"); p != "" {
		registryPath = p
	}
	downloadURL := registryPath + tarName
	_, _ = utils.ColorInfo.Printf("[Info] Downloading from: %s \n", downloadURL)

	dstDownloadDir := filepath.Join(gb.downloadsDir)
	_, _ = utils.ColorInfo.Printf("[Info] Downloading to: %s \n", dstDownloadDir)
	err := utils.DownloadWithProgress(downloadURL, tarName, dstDownloadDir)

	if err != nil {
		gb.cleanVersionDir(version)
		_, _ = utils.ColorInfo.Printf("[Info]: Downloading version failed: %s \n", err)
		_, _ = utils.ColorError.Printf("[Error]: Please check connectivity to url: %s\n", downloadURL)
		os.Exit(1)
	}

	srcTar := filepath.Join(gb.downloadsDir, tarName)
	dstDir := gb.getVersionDir(version)

	_, _ = utils.ColorInfo.Printf("[Info] Extracting from: %s \n", srcTar)
	_, _ = utils.ColorInfo.Printf("[Info] Extracting to: %s \n", dstDir)

	err = gb.ExtractTarGz(srcTar, dstDir)
	if err != nil {
		// clean up dir
		gb.cleanVersionDir(version)
		_, _ = utils.ColorInfo.Printf("[Info]: Untar failed: %s \n", err)
		_, _ = utils.ColorError.Printf("[Error]: Please check if version exists from url: %s\n", downloadURL)
		os.Exit(1)
	}
	_, _ = utils.ColorInfo.Printf("[Success] Untar to %s\n", gb.getVersionDir(version))
}

func (gb *GoBrew) ExtractTarGz(srcTar string, dstDir string) error {
	//#nosec G304
	file, err := os.Open(srcTar)
	if err != nil {
		return err
	}
	_, err = unpackit.Unpack(file, dstDir)
	if err != nil {
		return err
	}

	return nil
}

func (gb *GoBrew) changeSymblinkGoBin(version string) {

	goBinDst := filepath.Join(gb.versionsDir, version, "/go/bin")
	_ = os.RemoveAll(gb.currentBinDir)

	cmd := exec.Command("ln", "-snf", goBinDst, gb.currentBinDir)

	_, err := cmd.Output()
	if err != nil {
		_, _ = utils.ColorError.Printf("[Error]: symbolic link failed: %s\n", err)
		os.Exit(1)
	}

}

func (gb *GoBrew) changeSymblinkGo(version string) {
	_ = os.RemoveAll(gb.currentGoDir)
	versionGoDir := filepath.Join(gb.versionsDir, version, "go")
	cmd := exec.Command("ln", "-snf", versionGoDir, gb.currentGoDir)

	_, err := cmd.Output()
	if err != nil {
		_, _ = utils.ColorError.Printf("[Error]: symbolic link failed: %s\n", err)
		os.Exit(1)
	}
}
