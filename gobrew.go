package gobrew

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/c4milo/unpackit"
	"github.com/gookit/color"
	"github.com/kevincobain2000/gobrew/utils"
)

const (
	goBrewDir           string = ".gobrew"
	defaultRegistryPath string = "https://go.dev/dl/"
	goBrewDownloadUrl   string = "https://github.com/kevincobain2000/gobrew/releases/latest/download/"
	goBrewTagsApi       string = "https://raw.githubusercontent.com/kevincobain2000/gobrew/json/golang-tags.json"
)

// check GoBrew implement is Command interface
var _ Command = (*GoBrew)(nil)

// Command ...
type Command interface {
	ListVersions()
	ListRemoteVersions(print bool) map[string][]string
	CurrentVersion() string
	Uninstall(version string)
	Install(version string)
	Use(version string)
	Prune()
	Version(currentVersion string)
	Upgrade(currentVersion string)
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
}

var gb GoBrew

// NewGoBrew instance
func NewGoBrew() GoBrew {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.Getenv("HOME")
	}

	if os.Getenv("GOBREW_ROOT") != "" {
		homeDir = os.Getenv("GOBREW_ROOT")
	}

	return NewGoBrewDirectory(homeDir)
}

func NewGoBrewDirectory(homeDir string) GoBrew {
	gb.homeDir = homeDir

	gb.installDir = filepath.Join(gb.homeDir, goBrewDir)
	gb.versionsDir = filepath.Join(gb.installDir, "versions")
	gb.currentDir = filepath.Join(gb.installDir, "current")
	gb.currentBinDir = filepath.Join(gb.installDir, "current", "bin")
	gb.currentGoDir = filepath.Join(gb.installDir, "current", "go")
	gb.downloadsDir = filepath.Join(gb.installDir, "downloads")

	return gb
}

func (gb *GoBrew) Interactive(ask bool) {
	currentVersion := gb.CurrentVersion()
	currentMajorVersion := ExtractMajorVersion(currentVersion)

	latestVersion := gb.getLatestVersion()
	latestMajorVersion := ExtractMajorVersion(latestVersion)

	modVersion := gb.getModVersion()

	if modVersion == "" {
		modVersion = "None"
	}

	fmt.Println()

	if currentVersion == "" {
		currentVersion = "None"
		color.Warnln("GO Installed Version", ".......", currentVersion)
	} else {
		labels := []string{}
		if modVersion != "None" && currentMajorVersion != modVersion {
			labels = append(labels, "not same as go.mod")
		}
		if currentVersion != latestVersion {
			labels = append(labels, "not latest")
		}
		label := ""
		if len(labels) > 0 {
			label = "(" + strings.Join(labels, ", ") + ")"
		}
		if label != "" {
			label = " " + color.FgRed.Render(label)
		}
		color.Successln("GO Installed Version", ".......", currentVersion+label)
	}

	if latestMajorVersion != modVersion {
		label := " " + color.FgYellow.Render("(not latest)")
		color.Successln("GO go.mod Version", "   .......", modVersion+label)
	} else {
		color.Successln("GO go.mod Version", "   .......", modVersion)
	}

	color.Successln("GO Latest Version", "   .......", latestVersion)
	fmt.Println()

	if currentVersion == "None" {
		color.Warnln("GO is not installed.")
		c := true
		if ask {
			c = AskForConfirmation("Do you want to use latest GO version (" + latestVersion + ")?")
		}
		if c {
			gb.Install(latestVersion)
			gb.Use(latestVersion)
		}
		return
	}

	if currentMajorVersion != modVersion {
		color.Warnf("GO Installed Version (%s) and go.mod Version (%s) are different.\n", currentMajorVersion, modVersion)
		c := true
		if ask {
			c = AskForConfirmation("Do you want to use GO version same as go.mod version (" + modVersion + "@latest)?")
		}
		if c {
			gb.Install(modVersion + "@latest")
			gb.Use(modVersion + "@latest")
		}
		return
	}

	if currentVersion != latestVersion {
		color.Warnf("GO Installed Version (%s) and GO Latest Version (%s) are different.\n", currentVersion, latestVersion)
		c := true
		if ask {
			c = AskForConfirmation("Do you want to update GO to latest version (" + latestVersion + ")?")
		}
		if c {
			gb.Install(latestVersion)
			gb.Use(latestVersion)
		}
		return
	}
}

func (gb *GoBrew) getLatestVersion() string {
	getGolangVersions := gb.getGolangVersions()
	// loop through reverse and ignore beta and rc versions to get latest version
	for i := len(getGolangVersions) - 1; i >= 0; i-- {
		r := regexp.MustCompile("beta.*|rc.*")
		matches := r.FindAllString(getGolangVersions[i], -1)
		if len(matches) == 0 {
			return strings.ReplaceAll(getGolangVersions[i], "go", "")
		}
	}
	return ""
}
func (gb *GoBrew) getArch() string {
	return runtime.GOOS + "-" + runtime.GOARCH
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
func (gb *GoBrew) ListRemoteVersions(print bool) map[string][]string {
	color.Infoln("==> [Info] Fetching remote versions\n")
	tags := gb.getGolangVersions()

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
			r := regexp.MustCompile("beta.*|rc.*")
			matches := r.FindAllString(majorVersion, -1)
			if len(matches) == 1 {
				majorVersion = strings.Split(version, matches[0])[0]
			}
			if !isBlackListed(majorVersion) {
				groupedVersions[majorVersion] = append(groupedVersions[majorVersion], version)
			}
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
				color.Infop(versionParts[0])
			}
			gb.print("\t", print)
		} else {
			if print {
				color.Infop(lookupKey)
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

		maxPerLine = 0
		gb.print("\n\t", print)

		// print rc and beta versions in the end
		for _, rcVersion := range groupedVersions[lookupKey] {
			r := regexp.MustCompile("beta.*|rc.*")
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

func isBlackListed(version string) bool {
	blackListVersions := []string{"1.0", "1.1", "1.2", "1.3", "1.4"}
	for _, v := range blackListVersions {
		if version == v {
			return true
		}
	}
	return false
}

func (gb *GoBrew) print(message string, shouldPrint bool) {
	if shouldPrint {
		color.Infop(message)
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

	version := strings.ReplaceAll(fp, strings.Join([]string{"go", "bin"}, string(os.PathSeparator)), "")
	version = strings.ReplaceAll(version, gb.versionsDir, "")
	version = strings.ReplaceAll(version, string(os.PathSeparator), "")
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

func (gb *GoBrew) cleanVersionDir(version string) {
	_ = os.RemoveAll(gb.getVersionDir(version))
}

func (gb *GoBrew) cleanDownloadsDir() {
	_ = os.RemoveAll(gb.downloadsDir)
}

// Install the given version of go
func (gb *GoBrew) Install(version string) {
	if version == "" {
		color.Errorln("[Error] No version provided")
		os.Exit(1)
	}
	version = gb.judgeVersion(version)
	gb.mkDirs(version)
	if gb.existsVersion(version) {
		color.Infof("==> [Info] Version: %s exists\n", version)
		return
	}

	color.Infof("==> [Info] Downloading version: %s\n", version)
	gb.downloadAndExtract(version)
	gb.cleanDownloadsDir()
	color.Successf("==> [Success] Downloaded version: %s\n", version)
}

func (gb *GoBrew) judgeVersion(version string) string {
	judgedVersion := ""
	rcBetaOk := false
	reRcOrBeta := regexp.MustCompile("beta.*|rc.*")
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

	if version == "mod" {
		// get version by reading the mod file of Go
		judgedVersion = gb.getModVersion()
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
		if len(versionsSemantic) == 0 {
			return ""
		}

		// sort semantic versions
		sort.Sort(semver.Collection(versionsSemantic))
		// loop in reverse
		for i := len(versionsSemantic) - 1; i >= 0; i-- {
			judgedVersions := groupedVersions[versionsSemantic[i].Original()]
			// get last element
			if version == "dev-latest" {
				if len(judgedVersions) == 0 {
					return ""
				}
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

// read go.mod file and extract version
// Do not use go to get the version as go list -m -f '{{.GoVersion}}'
// Because go might not be installed
func (gb *GoBrew) getModVersion() string {
	modFilePath := filepath.Join("go.mod")
	modFile, err := os.Open(modFilePath)
	if err != nil {
		color.Errorln(err)
		os.Exit(1)
	}
	defer func(modFile *os.File) {
		_ = modFile.Close()
	}(modFile)

	scanner := bufio.NewScanner(modFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "go ") {
			return strings.TrimPrefix(line, "go ")
		}
	}

	if err = scanner.Err(); err != nil {
		color.Errorln(err)
		os.Exit(1)
	}
	return ""
}

// Use a version
func (gb *GoBrew) Use(version string) {
	version = gb.judgeVersion(version)
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

	fileExt := ""
	if runtime.GOOS == "windows" {
		fileExt = ".exe"
	}

	mkdirTemp, _ := os.MkdirTemp("", "gobrew")
	tmpFile := filepath.Join(mkdirTemp, "gobrew"+fileExt)
	url := goBrewDownloadUrl + "gobrew-" + gb.getArch() + fileExt
	utils.CheckError(
		utils.DownloadWithProgress(url, "gobrew"+fileExt, mkdirTemp),
		"[Error] Download GoBrew failed")

	source, err := os.Open(tmpFile)
	utils.CheckError(err, "[Error] Cannot open file")
	defer func(source *os.File) {
		_ = source.Close()
		utils.CheckError(os.Remove(source.Name()), "==> [Error] Cannot remove tmp file:")
	}(source)

	goBrewFile := filepath.Join(gb.installDir, "bin", "gobrew"+fileExt)
	if runtime.GOOS == "windows" {
		goBrewOldFile := goBrewFile + ".old"
		utils.CheckError(os.Rename(goBrewFile, goBrewOldFile), "==> [Error] Cannot rename binary file")
	} else {
		utils.CheckError(os.Remove(goBrewFile), "==> [Error] Cannot remove binary file")
	}
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

func (gb *GoBrew) mkDirs(version string) {
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
	tarName := "go" + version + "." + gb.getArch()

	if runtime.GOOS == "windows" {
		tarName = tarName + ".zip"
	} else {
		tarName = tarName + ".tar.gz"
	}

	registryPath := defaultRegistryPath
	if p := os.Getenv("GOBREW_REGISTRY"); p != "" {
		registryPath = p
	}
	downloadURL := registryPath + tarName
	color.Infoln("==> [Info] Downloading from:", downloadURL)

	dstDownloadDir := filepath.Join(gb.downloadsDir)
	color.Infoln("==> [Info] Downloading to:", dstDownloadDir)
	err := utils.DownloadWithProgress(downloadURL, tarName, dstDownloadDir)

	if err != nil {
		gb.cleanVersionDir(version)
		color.Infoln("==> [Info] Downloading version failed:", err)
		color.Errorln("==> [Error]: Please check connectivity to url:", downloadURL)
		os.Exit(1)
	}

	srcTar := filepath.Join(gb.downloadsDir, tarName)
	dstDir := gb.getVersionDir(version)

	color.Infoln("==> [Info] Extracting from:", srcTar)
	color.Infoln("==> [Info] Extracting to:", dstDir)

	err = gb.Extract(srcTar, dstDir)
	if err != nil {
		// clean up dir
		gb.cleanVersionDir(version)
		color.Infoln("==> [Info] Extract failed:", err)
		color.Errorln("==> [Error]: Please check if version exists from url:", downloadURL)
		os.Exit(1)
	}
	color.Infoln("[Success] Extract to", gb.getVersionDir(version))
}

func (gb *GoBrew) Extract(srcTar string, dstDir string) error {
	//#nosec G304
	file, err := os.Open(srcTar)
	if err != nil {
		return err
	}
	err = unpackit.Unpack(file, dstDir)
	if err != nil {
		return err
	}

	return nil
}

func (gb *GoBrew) changeSymblinkGoBin(version string) {
	goBinDst := filepath.Join(gb.versionsDir, version, "/go/bin")
	_ = os.RemoveAll(gb.currentBinDir)
	utils.CheckError(os.Symlink(goBinDst, gb.currentBinDir), "==> [Error]: symbolic link failed")
}

func (gb *GoBrew) changeSymblinkGo(version string) {
	_ = os.RemoveAll(gb.currentGoDir)
	versionGoDir := filepath.Join(gb.versionsDir, version, "go")
	utils.CheckError(os.Symlink(versionGoDir, gb.currentGoDir), "==> [Error]: symbolic link failed")
}

func (gb *GoBrew) getGobrewVersion() string {
	url := "https://api.github.com/repos/kevincobain2000/gobrew/releases/latest"
	data := doRequest(url)
	if len(data) == 0 {
		return ""
	}

	type Tag struct {
		TagName string `json:"tag_name"`
	}
	var tag Tag
	utils.CheckError(json.Unmarshal(data, &tag), "==> [Error]")

	return tag.TagName
}

func (gb *GoBrew) getGolangVersions() (result []string) {
	data := doRequest(goBrewTagsApi)
	if len(data) == 0 {
		return
	}

	type Tag struct {
		Ref string `json:"ref"`
	}
	var tags []Tag
	utils.CheckError(json.Unmarshal(data, &tags), "==> [Error]")

	for _, tag := range tags {
		t := strings.ReplaceAll(tag.Ref, "refs/tags/", "")
		if strings.HasPrefix(t, "go") {
			result = append(result, t)
		}
	}

	return
}

func doRequest(url string) (data []byte) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Errorln("==> [Error] Cannot create request:", err.Error())
		return
	}

	request.Header.Set("User-Agent", "gobrew")

	response, err := client.Do(request)
	utils.CheckError(err, "==> [Error] Cannot get response")

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(response.Body)

	if response.StatusCode == http.StatusTooManyRequests ||
		response.StatusCode == http.StatusForbidden {
		color.Errorln("==> [Error] Rate limit exhausted")
		os.Exit(1)
	}

	if response.StatusCode != http.StatusOK {
		color.Errorln("==> [Error] Cannot read response:", response.Status)
		os.Exit(1)
	}

	data, err = io.ReadAll(response.Body)
	utils.CheckError(err, "==> [Error] Cannot read response Body:")

	return
}
