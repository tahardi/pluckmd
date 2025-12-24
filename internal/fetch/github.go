package fetch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	GitHubHost    = "github.com"
	GitHubRawHost = "raw.githubusercontent.com"
	HTTPSPrefix   = "https://"
	BlobPath      = "/blob/"
	TreePath      = "/tree/"
)

var (
	ErrGitHubFetcher = errors.New("github fetcher")
	ErrBadURL        = errors.New("bad url")
)

type GitHubFetcher struct {
	client *http.Client
}

func NewGitHubFetcher() (*GitHubFetcher, error) {
	return NewGitHubFetcherWithClient(&http.Client{})
}

func NewGitHubFetcherWithClient(client *http.Client) (*GitHubFetcher, error) {
	return &GitHubFetcher{client: client}, nil
}

func (g *GitHubFetcher) Fetch(
	ctx context.Context,
	uri string,
) ([]byte, error) {
	url, err := URItoRawGitHubURL(uri)
	if err != nil {
		return nil, fmt.Errorf("%w: making url: %w", ErrGitHubFetcher, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: creating request: %w", ErrGitHubFetcher, err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: fetching: %w", ErrGitHubFetcher, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: reading body: %w", ErrGitHubFetcher, err)
	}
	return data, nil
}

func URItoRawGitHubURL(uri string) (string, error) {
	rawURL := strings.TrimSpace(uri)
	rawURL = strings.TrimSuffix(rawURL, "/")

	rawBase := HTTPSPrefix + GitHubRawHost
	if i := strings.Index(rawURL, GitHubHost); i != -1 {
		rawURL = rawURL[i+len(GitHubHost):]
	} else if j := strings.Index(rawURL, GitHubRawHost); j != -1 {
		rawURL = rawURL[j+len(GitHubRawHost):]
	} else {
		return "", fmt.Errorf("%w: '%s' is not a GitHub URL", ErrBadURL, rawURL)
	}
	rawURL = rawBase + rawURL

	rawURL = strings.Replace(rawURL, BlobPath, "/", 1)
	rawURL = strings.Replace(rawURL, TreePath, "/", 1)
	return rawURL, nil
}
