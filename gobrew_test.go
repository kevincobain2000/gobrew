package gobrew

import (
	"testing"

	"gotest.tools/assert"
)

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
			wantVersion: "1.18.4",
		},
		{
			version:     "1.18@dev-latest",
			wantVersion: "1.18.4",
		},
		// following 2 tests will fail upon new release of 1.19
		// at the time of this test, 1.19 is not released yet
		{
			version:     "latest",
			wantVersion: "1.18.4",
		},
		{
			version:     "dev-latest",
			wantVersion: "1.19rc2",
		},
	}
	for _, test := range tests {
		t.Run(test.version, func(t *testing.T) {
			gb := NewGoBrew()
			version := gb.judgeVersion(test.version)
			assert.Equal(t, test.wantVersion, version)

		})
	}
}
