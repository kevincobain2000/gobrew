package gobrew

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupGobrew(t *testing.T, ts *httptest.Server) GoBrew {
	tags, _ := url.JoinPath(ts.URL, "golang-tags.json")
	versionUrl, _ := url.JoinPath(ts.URL, "latest")
	config := Config{
		RootDir:           t.TempDir(),
		RegistryPathUrl:   ts.URL,
		GobrewDownloadUrl: ts.URL,
		GobrewTags:        tags,
		GobrewVersionsUrl: versionUrl,
	}
	gb := NewGoBrew(config)
	return gb
}

func TestInstallAndExistVersion(t *testing.T) {
	//t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)
	gb.Install("1.9")
	exists := gb.existsVersion("1.9")
	assert.Equal(t, true, exists)
	t.Log("test finished")
}

func TestUnInstallThenNotExistVersion(t *testing.T) {
	//t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)
	gb.Install("1.9")
	exists := gb.existsVersion("1.9")
	assert.Equal(t, true, exists)
	gb.Uninstall("1.9")
	exists = gb.existsVersion("1.9")
	assert.Equal(t, false, exists)
	t.Log("test finished")
}

func TestUpgrade(t *testing.T) {
	//t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)

	binaryDir := filepath.Join(gb.installDir, "bin")
	_ = os.MkdirAll(binaryDir, os.ModePerm)

	baseName := "gobrew" + fileExt
	binaryFile := filepath.Join(binaryDir, baseName)

	if oldFile, err := os.Create(binaryFile); err == nil {
		// on tests, we have to close the file to avoid an error on os.Rename
		_ = oldFile.Close()
	}

	gb.Upgrade("0.0.0")

	if _, err := os.Stat(binaryFile); err != nil {
		t.Errorf("updated executable does not exist")
	}
	t.Log("test finished")
}

func TestDoNotUpgradeLatestVersion(t *testing.T) {
	t.Skip("skipping test...needs to rewrite")
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)

	binaryDir := filepath.Join(gb.installDir, "bin")
	_ = os.MkdirAll(binaryDir, os.ModePerm)

	baseName := "gobrew" + fileExt
	binaryFile := filepath.Join(binaryDir, baseName)

	currentVersion := gb.getGobrewVersion()

	if currentVersion == "" {
		t.Skip("could not determine the current version")
	}

	gb.Upgrade(currentVersion[1:])

	if _, err := os.Stat(binaryFile); err == nil {
		t.Errorf("unexpected upgrade of latest version")
	}
	t.Log("test finished")
}

func TestInteractive(t *testing.T) {
	//t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)

	currentVersion := gb.CurrentVersion()
	latestVersion := gb.getLatestVersion()
	assert.Equal(t, "None", currentVersion)
	assert.NotEqual(t, currentVersion, latestVersion)

	gb.Interactive(false)

	currentVersion = gb.CurrentVersion()
	assert.Equal(t, currentVersion, latestVersion)

	gb.Install("1.16.5") // we know, it is not latest
	gb.Use("1.16.5")
	currentVersion = gb.CurrentVersion()
	assert.Equal(t, "1.16.5", currentVersion)
	assert.NotEqual(t, currentVersion, latestVersion)

	gb.Interactive(false)
	currentVersion = gb.CurrentVersion()
	assert.Equal(t, currentVersion, latestVersion)
	t.Log("test finished")
}

func TestPrune(t *testing.T) {
	//t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)
	gb.Install("1.20")
	gb.Install("1.19")
	gb.Use("1.19")
	gb.Prune()
	assert.Equal(t, false, gb.existsVersion("1.20"))
	assert.Equal(t, true, gb.existsVersion("1.19"))
	t.Log("test finished")
}

func TestGoBrew_CurrentVersion(t *testing.T) {
	//t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)
	assert.Equal(t, true, gb.CurrentVersion() == "None")
	gb.Install("1.19")
	gb.Use("1.19")
	assert.Equal(t, true, gb.CurrentVersion() == "1.19")
	t.Log("test finished")
}
