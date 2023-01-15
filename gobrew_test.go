package gobrew

import (
	"os"
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
	gb := NewGoBrew()
	err := gb.ListVersions()
	assert.NilError(t, err)
}
func TestExistVersion(t *testing.T) {
	gb := NewGoBrew()
	exists := gb.existsVersion("1.0") //ideally on tests nothing exists yet
	assert.Equal(t, false, exists)
}

func TestInstallAndExistVersion(t *testing.T) {
	gb := NewGoBrew()
	gb.Install("1.19")
	exists := gb.existsVersion("1.19")
	assert.Equal(t, true, exists)
}

func TestUnInstallThenNotExistVersion(t *testing.T) {
	gb := NewGoBrew()
	gb.Uninstall("1.19")
	exists := gb.existsVersion("1.19")
	assert.Equal(t, false, exists)
}
