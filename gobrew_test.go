package gobrew

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"gotest.tools/assert"
)

func TestNewGobrewHomeDirUsesUserHomeDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		t.FailNow()
	}

	gobrew := NewGoBrew()

	assert.Equal(t, homeDir, gobrew.homeDir)
}

func TestNewGobrewHomeDirDefaultsToHome(t *testing.T) {
	var envName string

	if runtime.GOOS == "windows" {
		envName = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		envName = "home"
	} else {
		envName = "HOME"
	}

	oldEnvValue := os.Getenv(envName)
	defer func() {
		os.Setenv(envName, oldEnvValue)
	}()

	os.Unsetenv(envName)

	gobrew := NewGoBrew()

	assert.Equal(t, os.Getenv("HOME"), gobrew.homeDir)
}

func TestNewGobrewHomeDirUsesGoBrewRoot(t *testing.T) {
	oldEnvValue := os.Getenv("GOBREW_ROOT")
	defer func() {
		os.Setenv("GOBREW_ROOT", oldEnvValue)
	}()

	os.Setenv("GOBREW_ROOT", "some_fancy_value")

	gobrew := NewGoBrew()

	assert.Equal(t, "some_fancy_value", gobrew.homeDir)
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
		// // following 2 tests fail upon new version release
		// // commenting out for now as the tool is stable
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
		t.Run(test.version, func(t *testing.T) {
			gb := NewGoBrew()
			version := gb.judgeVersion(test.version)
			assert.Equal(t, test.wantVersion, version)

		})
	}
}
func TestListVersions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gobrew-test-list-version-")
	if err != nil {
		t.Skip("could not create directory for gobrew update:", err)
		return
	}
	defer func() {
		os.RemoveAll(tempDir)
	}()

	gb := NewGoBrewDirectory(tempDir)
	err = gb.ListVersions()
	assert.NilError(t, err)
}

func TestExistVersion(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gobrew-test-exists-version-")
	if err != nil {
		t.Skip("could not create directory for gobrew update:", err)
		return
	}
	defer func() {
		os.RemoveAll(tempDir)
	}()

	gb := NewGoBrewDirectory(tempDir)
	exists := gb.existsVersion("1.19")
	assert.Equal(t, false, exists)
}

func TestInstallAndExistVersion(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "gobrew-test-install-uninstall")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Skip("could not create directory for gobrew update:", err)
		return
	}

	gb := NewGoBrewDirectory(tempDir)
	gb.Install("1.19")
	exists := gb.existsVersion("1.19")
	assert.Equal(t, true, exists)
}

func TestUnInstallThenNotExistVersion(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "gobrew-test-install-uninstall")
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		t.Skip("could not create directory for gobrew update:", err)
		return
	}
	defer func() {
		os.RemoveAll(tempDir)
	}()

	gb := NewGoBrewDirectory(tempDir)
	gb.Uninstall("1.19")
	exists := gb.existsVersion("1.19")
	assert.Equal(t, false, exists)
}

func TestUpgrade(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gobrew-test-upgrade-")
	if err != nil {
		t.Skip("could not create directory for gobrew update:", err)
		return
	}
	defer func() {
		os.RemoveAll(tempDir)
	}()

	gb := NewGoBrewDirectory(tempDir)

	binaryDir := filepath.Join(gb.installDir, "bin")
	_ = os.MkdirAll(binaryDir, os.ModePerm)
	goBrewFile := "gobrew"
	if runtime.GOOS == "windows" {
		goBrewFile = goBrewFile + ".exe"
	}

	gb.Upgrade("0.0.0")

	if _, err := os.Stat(filepath.Join(binaryDir, goBrewFile)); err != nil {
		t.Errorf("could not find updated executable")
	}
}
