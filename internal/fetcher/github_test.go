package fetcher_test

import (
	"testing"

	"github.com/tahardi/pluckmd/internal/fetcher"
)

func TestGitHubFetcher_Fetch(t *testing.T) {
	// TODO: Happy path fetch github.go from this repo
	// and compare to an embeded version!
}

func TestURItoRawGitHubURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Standard blob URL",
			input:    "https://github.com/user/repo/blob/main/file.go",
			expected: "https://raw.githubusercontent.com/user/repo/main/file.go",
		},
		{
			name:     "Standard tree URL",
			input:    "https://github.com/user/repo/tree/v1.0.0/README.md",
			expected: "https://raw.githubusercontent.com/user/repo/v1.0.0/README.md",
		},
		{
			name:     "No scheme",
			input:    "github.com/user/repo/blob/main/file.go",
			expected: "https://raw.githubusercontent.com/user/repo/main/file.go",
		},
		{
			name:     "Insecure HTTP scheme",
			input:    "http://github.com/user/repo/blob/main/file.go",
			expected: "https://raw.githubusercontent.com/user/repo/main/file.go",
		},
		{
			name:     "Already a raw URL",
			input:    "https://raw.githubusercontent.com/user/repo/main/file.go",
			expected: "https://raw.githubusercontent.com/user/repo/main/file.go",
		},
		{
			name:     "URL with leading/trailing whitespace",
			input:    "  https://github.com/user/repo/blob/main/file.go  ",
			expected: "https://raw.githubusercontent.com/user/repo/main/file.go",
		},
		{
			name:     "URL with leading jibberish scheme",
			input:    "jibberishscheme://github.com/user/repo/blob/main/file.go  ",
			expected: "https://raw.githubusercontent.com/user/repo/main/file.go",
		},
		{
			name:     "Trailing slash",
			input:    "https://github.com/user/repo/blob/main/file.go/",
			expected: "https://raw.githubusercontent.com/user/repo/main/file.go",
		},
		{
			name:     "Deeply nested path",
			input:    "https://github.com/user/repo/blob/main/pkg/sub/module/file.go",
			expected: "https://raw.githubusercontent.com/user/repo/main/pkg/sub/module/file.go",
		},
		{
			name:    "Invalid host",
			input:   "https://gitlab.com/user/repo/blob/main/file.go",
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetcher.URItoRawGitHubURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("URItoRawGitHubURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("URItoRawGitHubURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}
