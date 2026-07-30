package inspect

import "testing"

func TestAbsolutePathPatterns(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "unix home path",
			input: "/home/kato/project/main.go",
			want:  true,
		},
		{
			name:  "mac users path",
			input: "/Users/kato/project/main.go",
			want:  true,
		},
		{
			name:  "windows users path",
			input: `C:\Users\kato\work\main.go`,
			want:  true,
		},
		{
			name:  "known system mime path",
			input: "/usr/local/share/mime/globs2",
			want:  false,
		},
		{
			name:  "short backslash noise",
			input: `\\b\c\g\h\i\m\p\t\u\z\`,
			want:  false,
		},
		{
			name:  "unc path",
			input: `\\server01\share01\project\main.go`,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := false
			for _, re := range absolutePathPatterns {
				got = got || re.MatchString(tt.input)
			}
			if got != tt.want {
				t.Fatalf("match = %v, want %v", got, tt.want)
			}
		})
	}
}
