package fetch

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrLocalFetcher = errors.New("local fetcher")
)

type LocalFletcher struct {
	baseDir string
}

func NewLocalFetcher() (*LocalFletcher, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("%w: getting working directory: %w", ErrLocalFetcher, err)
	}
	return NewLocalFetcherWithBaseDir(cwd)
}

func NewLocalFetcherWithBaseDir(baseDir string) (*LocalFletcher, error) {
	return &LocalFletcher{baseDir: baseDir}, nil
}

func (l *LocalFletcher) Fetch(
	_ context.Context,
	uri string,
) ([]byte, error) {
	if !filepath.IsAbs(uri) {
		uri = filepath.Join(l.baseDir, uri)
	}

	data, err := os.ReadFile(uri)
	if err != nil {
		return nil, fmt.Errorf("%w: reading file '%s': %w", ErrLocalFetcher, uri, err)
	}
	return data, nil
}
