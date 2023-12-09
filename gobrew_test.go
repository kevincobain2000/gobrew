package gobrew

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallAndExistVersion(t *testing.T) {
	t.Parallel()
	gb := NewGoBrew(t.TempDir())
	gb.Install("1.19")
	exists := gb.existsVersion("1.19")
	assert.Equal(t, true, exists)
	t.Log("test finished")
}

func TestUnInstallThenNotExistVersion(t *testing.T) {
	t.Parallel()
	gb := NewGoBrew(t.TempDir())
	gb.Uninstall("1.19")
	exists := gb.existsVersion("1.19")
	assert.Equal(t, false, exists)
	t.Log("test finished")
}

func TestUpgrade(t *testing.T) {
	t.Parallel()
	gb := NewGoBrew(t.TempDir())

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
	gb := NewGoBrew(t.TempDir())

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
	t.Parallel()
	gb := NewGoBrew(t.TempDir())
	currentVersion := gb.CurrentVersion()
	latestVersion := gb.getLatestVersion()
	// modVersion := gb.getModVersion()
	assert.Equal(t, "None", currentVersion)
	assert.NotEqual(t, currentVersion, latestVersion)

	gb.Interactive(false)

	currentVersion = gb.CurrentVersion()
	// remove string private from currentVersion (for macOS) due to /private/var symlink issue
	currentVersion = strings.Replace(currentVersion, "private", "", -1)
	assert.Equal(t, currentVersion, latestVersion)

	gb.Install("1.16.5") // we know, it is not latest
	gb.Use("1.16.5")
	currentVersion = gb.CurrentVersion()
	currentVersion = strings.Replace(currentVersion, "private", "", -1)
	assert.Equal(t, "1.16.5", currentVersion)
	assert.NotEqual(t, currentVersion, latestVersion)

	gb.Interactive(false)
	currentVersion = gb.CurrentVersion()
	currentVersion = strings.Replace(currentVersion, "private", "", -1)
	assert.Equal(t, currentVersion, latestVersion)
	t.Log("test finished")
}

func TestPrune(t *testing.T) {
	t.Parallel()
	gb := NewGoBrew(t.TempDir())
	gb.Install("1.20")
	gb.Install("1.19")
	gb.Use("1.19")
	gb.Prune()
	assert.Equal(t, false, gb.existsVersion("1.20"))
	assert.Equal(t, true, gb.existsVersion("1.19"))
	t.Log("test finished")
}

func TestGoBrew_CurrentVersion(t *testing.T) {
	t.Parallel()
	gb := NewGoBrew(t.TempDir())
	assert.Equal(t, true, gb.CurrentVersion() == "None")
	gb.Install("1.19")
	gb.Use("1.19")
	assert.Equal(t, true, gb.CurrentVersion() == "1.19")
	t.Log("test finished")
}
