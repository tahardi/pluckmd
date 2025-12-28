package run

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tahardi/pluckmd/internal/cache"
	"github.com/tahardi/pluckmd/internal/fetch"
	"github.com/tahardi/pluckmd/internal/pluck"
	"github.com/tahardi/pluckmd/internal/process"
)

const (
	DefaultPermissions = 0644
	MarkdownExt        = ".md"
)

var (
	ErrRunner = errors.New("runner")
)

type Runner struct {
	processor *process.Processor
}

func NewRunner() (*Runner, error) {
	cacher, err := cache.NewRAMCacher()
	if err != nil {
		return nil, err
	}

	ghFetcher, err := fetch.NewGitHubFetcher()
	if err != nil {
		return nil, err
	}

	lFetcher, err := fetch.NewLocalFetcher()
	if err != nil {
		return nil, err
	}
	fetchers := []fetch.Fetcher{ghFetcher, lFetcher}

	plucker, err := pluck.NewBlockyPlucker()
	if err != nil {
		return nil, err
	}
	return NewRunnerWithProcessor(process.NewProcessor(cacher, fetchers, plucker))
}

func NewRunnerWithProcessor(processor *process.Processor) (*Runner, error) {
	return &Runner{processor: processor}, nil
}

func (r *Runner) Run(ctx context.Context, dir string) error {
	files, err := ListMarkdownFiles(dir)
	if err != nil {
		return fmt.Errorf("%w: listing markdown files: %w", ErrRunner, err)
	}

	for _, file := range files {
		bytes, readErr := os.ReadFile(file)
		if readErr != nil {
			return fmt.Errorf("%w: reading file: %w", ErrRunner, readErr)
		}

		processed, procErr := r.processor.ProcessMarkdown(ctx, bytes)
		if procErr != nil {
			return fmt.Errorf("%w: processing file: %w", ErrRunner, procErr)
		}

		// #nosec G306
		writeErr := os.WriteFile(file, processed, DefaultPermissions)
		if writeErr != nil {
			return fmt.Errorf("%w: writing file: %w", ErrRunner, writeErr)
		}
	}
	return nil
}

func ListMarkdownFiles(dir string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == MarkdownExt {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
