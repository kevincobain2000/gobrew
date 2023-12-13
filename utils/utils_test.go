package utils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDownloadWithProgress(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.FileServer(http.Dir("../testdata")))
	defer ts.Close()
	path, _ := url.JoinPath(ts.URL, "go1.9.darwin-arm64.tar.gz")
	if err := DownloadWithProgress(path, "go1.9.darwin-arm64.tar.gz", t.TempDir()); (err != nil) != false {
		t.Errorf("DownloadWithProgress() error = %v, wantErr %v", err, false)
	}
	t.Log("test finished")
}
