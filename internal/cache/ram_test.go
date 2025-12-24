
package cache_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/pluckmd/internal/cache"
)

func TestRAMCacher(t *testing.T) {
	ctx := context.Background()

	t.Run("happy path", func(t *testing.T) {
		// given
		cacher, err := cache.NewRAMCacher()
		require.NoError(t, err)
		uri := "https://example.com/data"
		content := []byte("hello world")

		// when
		err = cacher.Store(ctx, uri, content)
		require.NoError(t, err)
		retrieved, err := cacher.Retrieve(ctx, uri)

		// then
		require.NoError(t, err)
		assert.Equal(t, content, retrieved)
	})

	t.Run("retrieve uri not found", func(t *testing.T) {
		// given
		cacher, _ := cache.NewRAMCacher()
		uri := "non-existent"

		// when
		_, err := cacher.Retrieve(ctx, uri)

		// then
		require.Error(t, err)
		assert.True(t, errors.Is(err, cache.ErrRAMCacher))
		assert.True(t, errors.Is(err, cache.ErrURINotFound))
		assert.Contains(t, err.Error(), uri)
	})

	t.Run("close clears the vault", func(t *testing.T) {
		// given
		cacher, _ := cache.NewRAMCacher()
		uri := "key"
		data := []byte("value")
		err := cacher.Store(ctx, uri, data)
		require.NoError(t, err)

		got, err := cacher.Retrieve(ctx, uri)
		require.NoError(t, err)
		assert.Equal(t, data, got)

		// when
		err = cacher.Close()
		require.NoError(t, err)
		_, retrieveErr := cacher.Retrieve(ctx, uri)

		// then
		require.Error(t, retrieveErr)
		assert.True(t, errors.Is(retrieveErr, cache.ErrURINotFound))
	})

	t.Run("initialize with existing vault", func(t *testing.T) {
		// given
		vault := map[string][]byte{
			"pre-existing": []byte("initial"),
		}

		// when
		cacher, err := cache.NewRAMCacherWithVault(vault)
		require.NoError(t, err)
		data, retrieveErr := cacher.Retrieve(ctx, "pre-existing")

		// then
		require.NoError(t, retrieveErr)
		assert.Equal(t, []byte("initial"), data)
	})
}
