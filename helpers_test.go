package gobrew

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJudgeVersion(t *testing.T) {
	t.Parallel()
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
		{
			version:     "go1.18",
			wantVersion: "None",
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
			t.Parallel()
			ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
			defer ts.Close()
			gb := setupGobrew(t, ts)
			version := gb.judgeVersion(test.version)
			assert.Equal(t, test.wantVersion, version)
		})
	}
	t.Log("test finished")
}

func TestListVersions(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)

	gb.ListVersions()
	t.Log("test finished")
}

func TestExistVersion(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)

	exists := gb.existsVersion("1.19")

	assert.Equal(t, false, exists)
	t.Log("test finished")
}

func TestExtractMajorVersion(t *testing.T) {
	t.Parallel()
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "extracts major version",
			args: args{
				version: "1.16.5",
			},
			want: "1.16",
		},
		{
			name: "extracts major version if already major version",
			args: args{
				version: "1.16",
			},
			want: "1.16",
		},
		{
			name: "extracts major version with extra parts",
			args: args{
				version: "",
			},
			want: "",
		},
		{
			name: "extracts major version rc",
			args: args{
				version: "1.21rc1",
			},
			want: "1.21",
		},
		{
			name: "extracts major version beta",
			args: args{
				version: "1.21beta1",
			},
			want: "1.21",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := extractMajorVersion(tt.args.version); got != tt.want {
				t.Errorf("ExtractMajorVersion() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Log("test finished")
}

func TestGoBrew_extract(t *testing.T) {
	t.Parallel()
	type args struct {
		srcTar string
		dstDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "go1.9.darwin-arm64.tar.gz",
			args: args{
				srcTar: "testdata/go1.9.darwin-arm64.tar.gz",
				dstDir: "tmp",
			},
			wantErr: false,
		},
		{
			name: "dont.tar.gz",
			args: args{
				srcTar: "testdata/dont.tar.gz",
				dstDir: "tmp",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
			defer ts.Close()
			gb := setupGobrew(t, ts)
			if err := gb.extract(tt.args.srcTar, filepath.Join(t.TempDir(), tt.args.dstDir)); (err != nil) != tt.wantErr {
				t.Errorf("GoBrew.extract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Log("test finished")
}

func Test_doRequest(t *testing.T) {
	t.Parallel()
	type args struct {
		url string
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
	}{
		{
			name: "test.txt",
			args: args{
				url: "test.txt",
			},
			wantData: []byte("test"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
			defer ts.Close()
			urlGet, _ := url.JoinPath(ts.URL, tt.args.url)
			if gotData := doRequest(urlGet); !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("doRequest() = %s, want %s", gotData, tt.wantData)
			}
		})
	}
	t.Log("test finished")
}

func TestGoBrew_downloadAndExtract(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	gb := setupGobrew(t, ts)
	gb.mkDirs("1.9")
	gb.downloadAndExtract("1.9")
	t.Log("test finished")
}
