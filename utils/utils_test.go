package utils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDownloadWithProgress(t *testing.T) {
	t.Skip()
	type args struct {
		name       string
		destFolder string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				name:       "go1.9.darwin-arm64.tar.gz",
				destFolder: t.TempDir(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
			defer ts.Close()
			url, _ := url.JoinPath(ts.URL, tt.args.name)
			if err := DownloadWithProgress(url, tt.args.name, tt.args.destFolder); (err != nil) != tt.wantErr {
				t.Errorf("DownloadWithProgress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
