package process_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/pluckmd/internal/cache"
	"github.com/tahardi/pluckmd/internal/fetch"
	"github.com/tahardi/pluckmd/internal/mocks"
	"github.com/tahardi/pluckmd/internal/pluck"
	"github.com/tahardi/pluckmd/internal/process"
)

const (
	processorCodeSnippetURI = "./processor.go.type.Processor"
	processorSourceCodeURI  = "./processor.go"
)

//go:embed testdata/processed.md
var processedMD []byte

//go:embed testdata/unprocessed.md
var unprocessedMD []byte

func TestProcessor_ProcessMarkdown(t *testing.T) {
	t.Run("happy path - unprocessed", func(t *testing.T) {
		// given
		cacher, err := cache.NewRAMCacher()
		require.NoError(t, err)

		fetcher, err := fetch.NewLocalFetcher()
		require.NoError(t, err)

		goPlucker, err := pluck.NewGoPlucker()
		require.NoError(t, err)
		yamlPlucker, err := pluck.NewYAMLPlucker()
		require.NoError(t, err)

		fetchers := []fetch.Fetcher{fetcher}
		pluckers := map[pluck.Lang]pluck.Plucker{
			pluck.Go:   goPlucker,
			pluck.YAML: yamlPlucker,
		}

		processor := process.NewProcessor(cacher, fetchers, pluckers)

		// when
		got, err := processor.ProcessMarkdown(context.Background(), unprocessedMD)

		// then
		require.NoError(t, err)
		require.Equal(t, processedMD, got)
	})

	t.Run("happy path - already processed", func(t *testing.T) {
		// given
		cacher, err := cache.NewRAMCacher()
		require.NoError(t, err)

		fetcher, err := fetch.NewLocalFetcher()
		require.NoError(t, err)

		goPlucker, err := pluck.NewGoPlucker()
		require.NoError(t, err)
		yamlPlucker, err := pluck.NewYAMLPlucker()
		require.NoError(t, err)

		fetchers := []fetch.Fetcher{fetcher}
		pluckers := map[pluck.Lang]pluck.Plucker{
			pluck.Go:   goPlucker,
			pluck.YAML: yamlPlucker,
		}

		processor := process.NewProcessor(cacher, fetchers, pluckers)

		// when
		got, err := processor.ProcessMarkdown(context.Background(), processedMD)

		// then
		require.NoError(t, err)
		require.Equal(t, processedMD, got)
	})

	t.Run("error - cacher", func(t *testing.T) {
		// given
		ctx := context.Background()
		cacher := mocks.NewCacher(t)
		cacher.On("Retrieve", ctx, processorCodeSnippetURI).Return(nil, assert.AnError)

		fetcher, err := fetch.NewLocalFetcher()
		require.NoError(t, err)

		goPlucker, err := pluck.NewGoPlucker()
		require.NoError(t, err)
		yamlPlucker, err := pluck.NewYAMLPlucker()
		require.NoError(t, err)

		fetchers := []fetch.Fetcher{fetcher}
		pluckers := map[pluck.Lang]pluck.Plucker{
			pluck.Go:   goPlucker,
			pluck.YAML: yamlPlucker,
		}

		processor := process.NewProcessor(cacher, fetchers, pluckers)

		// when
		_, err = processor.ProcessMarkdown(context.Background(), unprocessedMD)

		// then
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("error - fetcher", func(t *testing.T) {
		// given
		ctx := context.Background()
		cacher, err := cache.NewRAMCacher()
		require.NoError(t, err)

		fetcher := mocks.NewFetcher(t)
		fetcher.On("Fetch", ctx, processorSourceCodeURI).Return(nil, assert.AnError)

		goPlucker, err := pluck.NewGoPlucker()
		require.NoError(t, err)
		yamlPlucker, err := pluck.NewYAMLPlucker()
		require.NoError(t, err)

		fetchers := []fetch.Fetcher{fetcher}
		pluckers := map[pluck.Lang]pluck.Plucker{
			pluck.Go:   goPlucker,
			pluck.YAML: yamlPlucker,
		}

		processor := process.NewProcessor(cacher, fetchers, pluckers)

		// when
		_, err = processor.ProcessMarkdown(context.Background(), unprocessedMD)

		// then
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("error - plucker", func(t *testing.T) {
		// given
		ctx := context.Background()
		cacher, err := cache.NewRAMCacher()
		require.NoError(t, err)

		fetcher, err := fetch.NewLocalFetcher()
		require.NoError(t, err)

		goPlucker := mocks.NewPlucker(t)
		goPlucker.On("Pluck", ctx, mock.Anything, mock.Anything, pluck.Type).Return("", assert.AnError)

		yamlPlucker, err := pluck.NewYAMLPlucker()
		require.NoError(t, err)

		fetchers := []fetch.Fetcher{fetcher}
		pluckers := map[pluck.Lang]pluck.Plucker{
			pluck.Go:   goPlucker,
			pluck.YAML: yamlPlucker,
		}

		processor := process.NewProcessor(cacher, fetchers, pluckers)

		// when
		_, err = processor.ProcessMarkdown(context.Background(), unprocessedMD)

		// then
		assert.ErrorIs(t, err, assert.AnError)
	})
}
