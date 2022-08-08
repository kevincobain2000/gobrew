package gobrew

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
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
	goBrewDownloadUrl   string = "https://github.com/kevincobain2000/gobrew/releases/latest/download/"
)

// Command ...
type Command interface {
	ListVersions()
	ListRemoteVersions()
	CurrentVersion() string
	Uninstall(version string)
	Install(version string)
	Use(version string)
	Upgrade()
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
	getLatestVersion() string
	getGithubTags(repo string) (result []string)
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
	utils.CheckError(err, "[Error]: List versions failed")
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
		// 1.8.0 -> 1.8
		reMajorVersion, _ := regexp.Compile("[0-9]+.[0-9]+.0")
		if reMajorVersion.MatchString(version) {
			version = strings.Split(version, ".")[0] + "." + strings.Split(version, ".")[1]
		}

		if version == cv {
			version = cv + "*"
			utils.Successln(version)
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
				utils.Successln(rcVersion)
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
	tags := gb.getGithubTags("golang/go")

	var versions []string
	for _, tag := range tags {
		versions = append(versions, strings.ReplaceAll(tag, "go", ""))
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
		if v, err := semver.NewVersion(r); err == nil {
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
				utils.Major(versionParts[0])
			}
			gb.print("\t", print)
		} else {
			if print {
				utils.Major(lookupKey)
			}
			gb.print("\t", print)
		}

		groupedVersionsSemantic := make([]*semver.Version, 0)
		for _, r := range groupedVersions[lookupKey] {
			if v, err := semver.NewVersion(r); err == nil {
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
		utils.Errorf("[Error] Version: %s you are trying to remove is your current version. Please use a different version first before uninstalling the current version\n", version)
		os.Exit(1)
	}
	if !gb.existsVersion(version) {
		utils.Errorf("[Error] Version: %s you are trying to remove is not installed\n", version)
		os.Exit(1)
	}
	gb.cleanVersionDir(version)
	utils.Successf("[Success] Version: %s uninstalled\n", version)
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
		utils.Infof("[Info] Version: %s exists \n", version)
		return
	}

	utils.Infof("[Info] Downloading version: %s \n", version)
	gb.downloadAndExtract(version)
	gb.cleanDownloadsDir()
	utils.Successf("[Success] Downloaded version: %s\n", version)
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
			if v, err := semver.NewVersion(r); err == nil {
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
		utils.Infof("[Info] Version: %s is already your current version \n", version)
		return
	}
	utils.Infof("[Info] Changing go version to: %s \n", version)
	gb.changeSymblinkGoBin(version)
	gb.changeSymblinkGo(version)
	utils.Successf("[Success] Changed go version to: %s\n", version)
}

// Upgrade of GoBrew
func (gb *GoBrew) Upgrade(currentVersion string) {
	if "v"+currentVersion == gb.getLatestVersion() {
		utils.Infoln("[INFO] your version is already newest")
		return
	}

	mkdirTemp, _ := os.MkdirTemp("", "gobrew")
	tmpFile := filepath.Join(mkdirTemp, "gobrew")
	url := goBrewDownloadUrl + "gobrew-" + gb.getArch()
	if err := utils.DownloadWithProgress(url, "gobrew", mkdirTemp); err != nil {
		utils.Errorln("[Error] Download GoBrew failed:", err)
		return
	}

	source, err := os.Open(tmpFile)
	if err != nil {
		utils.Errorln("[Error] Cannot open file", err)
		return
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)

	goBrewFile := filepath.Join(gb.installDir, "/bin/gobrew")
	destination, err := os.Create(goBrewFile)
	if err != nil {
		utils.Errorf("[Error] Cannot open file: %s", err)
		return
	}
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)

	if _, err = io.Copy(destination, source); err != nil {
		utils.Errorf("[Error] Cannot copy file: %s", err)
		return
	}

	if err = os.Chmod(goBrewFile, 0755); err != nil {
		utils.Errorf("[Error] Cannot set file as executable: %s", err)
		return
	}

	if err = os.Remove(tmpFile); err != nil {
		utils.Errorf("[Error] Cannot remove tmp file: %s", err)
		return
	}

	utils.Infoln("Upgrade successful")
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
	utils.Infof("[Info] Downloading from: %s \n", downloadURL)

	dstDownloadDir := filepath.Join(gb.downloadsDir)
	utils.Infof("[Info] Downloading to: %s \n", dstDownloadDir)
	err := utils.DownloadWithProgress(downloadURL, tarName, dstDownloadDir)

	if err != nil {
		gb.cleanVersionDir(version)
		utils.Infof("[Info]: Downloading version failed: %s \n", err)
		utils.Errorf("[Error]: Please check connectivity to url: %s\n", downloadURL)
		os.Exit(1)
	}

	srcTar := filepath.Join(gb.downloadsDir, tarName)
	dstDir := gb.getVersionDir(version)

	utils.Infof("[Info] Extracting from: %s \n", srcTar)
	utils.Infof("[Info] Extracting to: %s \n", dstDir)

	err = gb.ExtractTarGz(srcTar, dstDir)
	if err != nil {
		// clean up dir
		gb.cleanVersionDir(version)
		utils.Infof("[Info]: Untar failed: %s \n", err)
		utils.Errorf("[Error]: Please check if version exists from url: %s\n", downloadURL)
		os.Exit(1)
	}
	utils.Infof("[Success] Untar to %s\n", gb.getVersionDir(version))
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

	if err := os.Symlink(goBinDst, gb.currentBinDir); err != nil {
		utils.Errorf("[Error]: symbolic link failed: %s\n", err)
		os.Exit(1)
	}
}

func (gb *GoBrew) changeSymblinkGo(version string) {
	_ = os.RemoveAll(gb.currentGoDir)
	versionGoDir := filepath.Join(gb.versionsDir, version, "go")

	if err := os.Symlink(versionGoDir, gb.currentGoDir); err != nil {
		utils.Errorf("[Error]: symbolic link failed: %s\n", err)
		os.Exit(1)
	}
}

func (gb *GoBrew) getLatestVersion() string {
	tags := gb.getGithubTags("kevincobain2000/gobrew")

	return tags[len(tags)-1]
}

func (gb *GoBrew) getGithubTags(repo string) (result []string) {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/git/refs/tags", repo), nil)
	if err != nil {
		utils.Errorf("[Error] Cannot create request: %s", err)
		return
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		utils.Errorf("[Error] Cannot get response: %s", err)
		return
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		utils.Errorf("[Error] Cannot read response: %s", err)
		return
	}

	type Tag struct {
		Ref string
	}
	var tags []Tag

	if err := json.Unmarshal(data, &tags); err != nil {
		utils.Errorf("[Error] Cannot unmarshal data: %s", err)
	}

	for _, tag := range tags {
		t := strings.ReplaceAll(tag.Ref, "refs/tags/", "")
		if strings.HasPrefix(t, "v") || strings.HasPrefix(t, "go") {
			result = append(result, t)
		}
	}

	return result
}
