package gobrew

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGobrewHomeDirUsesUserHomeDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		t.FailNow()
	}

	gobrew := NewGoBrew()

	assert.Equal(t, homeDir, gobrew.homeDir)
	t.Log("test finished")
}

func TestNewGobrewHomeDirDefaultsToHome(t *testing.T) {
	var envName string

	switch runtime.GOOS {
	case "windows":
		envName = "USERPROFILE"
	case "plan9":
		envName = "home"
	default:
		envName = "HOME"
	}

	oldEnvValue := os.Getenv(envName)
	defer func() {
		_ = os.Setenv(envName, oldEnvValue)
	}()

	_ = os.Unsetenv(envName)

	gobrew := NewGoBrew()

	assert.Equal(t, os.Getenv("HOME"), gobrew.homeDir)
	t.Log("test finished")
}

func TestNewGobrewHomeDirUsesGoBrewRoot(t *testing.T) {
	oldEnvValue := os.Getenv("GOBREW_ROOT")
	defer func() {
		_ = os.Setenv("GOBREW_ROOT", oldEnvValue)
	}()

	_ = os.Setenv("GOBREW_ROOT", "some_fancy_value")

	gobrew := NewGoBrew()

	assert.Equal(t, "some_fancy_value", gobrew.homeDir)
	t.Log("test finished")
}

func TestJudgeVersion(t *testing.T) {
	tests := []struct {
		version     string
		wantVersion string
		wantError   error
	}{
		{
			version:     "1.8",
			wantVersion: "1.8",
		},
		{
			version:     "1.8.2",
			wantVersion: "1.8.2",
		},
		{
			version:     "1.18beta1",
			wantVersion: "1.18beta1",
		},
		{
			version:     "1.18rc1",
			wantVersion: "1.18rc1",
		},
		{
			version:     "1.18@latest",
			wantVersion: "1.18.10",
		},
		{
			version:     "1.18@dev-latest",
			wantVersion: "1.18.10",
		},
		// following 2 tests fail upon new version release
		// commenting out for now as the tool is stable
		// {
		// 	version:     "latest",
		// 	wantVersion: "1.19.1",
		// },
		// {
		// 	version:     "dev-latest",
		// 	wantVersion: "1.19.1",
		// },
	}
	for _, test := range tests {
		test := test
		t.Run(test.version, func(t *testing.T) {
			gb := NewGoBrew()
			version := gb.judgeVersion(test.version)
			assert.Equal(t, test.wantVersion, version)

		})
	}
	t.Log("test finished")
}

func TestListVersions(t *testing.T) {
	tempDir := t.TempDir()
	gb := NewGoBrewDirectory(tempDir)

	gb.ListVersions()
	t.Log("test finished")
}

func TestExistVersion(t *testing.T) {
	tempDir := t.TempDir()
	gb := NewGoBrewDirectory(tempDir)

	exists := gb.existsVersion("1.19")

	assert.Equal(t, false, exists)
	t.Log("test finished")
}

func TestInstallAndExistVersion(t *testing.T) {
	tempDir := t.TempDir()
	gb := NewGoBrewDirectory(tempDir)
	gb.Install("1.19")
	exists := gb.existsVersion("1.19")
	assert.Equal(t, true, exists)
	t.Log("test finished")
}

func TestUnInstallThenNotExistVersion(t *testing.T) {
	tempDir := t.TempDir()
	gb := NewGoBrewDirectory(tempDir)
	gb.Uninstall("1.19")
	exists := gb.existsVersion("1.19")
	assert.Equal(t, false, exists)
	t.Log("test finished")
}

func TestUpgrade(t *testing.T) {
	tempDir := t.TempDir()

	gb := NewGoBrewDirectory(tempDir)

	binaryDir := filepath.Join(gb.installDir, "bin")
	_ = os.MkdirAll(binaryDir, os.ModePerm)

	baseName := "gobrew" + fileExt
	binaryFile := filepath.Join(binaryDir, baseName)

	if oldFile, err := os.Create(binaryFile); err == nil {
		// on tests we have to close the file to avoid an error on os.Rename
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
	tempDir := t.TempDir()

	gb := NewGoBrewDirectory(tempDir)

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
	tempDir := t.TempDir()

	gb := NewGoBrewDirectory(tempDir)
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
	tempDir := t.TempDir()
	gb := NewGoBrewDirectory(tempDir)
	gb.Install("1.20")
	gb.Install("1.19")
	gb.Use("1.19")
	gb.Prune()
	assert.Equal(t, false, gb.existsVersion("1.20"))
	assert.Equal(t, true, gb.existsVersion("1.19"))
	t.Log("test finished")
}

func TestGoBrew_CurrentVersion(t *testing.T) {
	tempDir := t.TempDir()
	gb := NewGoBrewDirectory(tempDir)
	assert.Equal(t, true, gb.CurrentVersion() == "None")
	gb.Install("1.19")
	gb.Use("1.19")
	assert.Equal(t, true, gb.CurrentVersion() == "1.19")
	t.Log("test finished")
}
