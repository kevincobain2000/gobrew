package gobrew

import (
	"testing"
)

func TestExtractMajorVersion(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractMajorVersion(tt.args.version); got != tt.want {
				t.Errorf("ExtractMajorVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
