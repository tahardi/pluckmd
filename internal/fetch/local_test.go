package fetch_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/pluckmd/internal/fetch"
)

const localFetcherURI = "./local.go"

//go:embed local.go
var localGo []byte

func TestLocalFetcher_Fetch(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		ctx := context.Background()
		want := localGo
		uri := localFetcherURI
		fetcher, err := fetch.NewLocalFetcher()
		require.NoError(t, err)

		// when
		got, err := fetcher.Fetch(ctx, uri)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
