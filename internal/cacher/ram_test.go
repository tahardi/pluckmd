
package cacher_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/pluckmd/internal/cacher"
)

func TestRAMCacher(t *testing.T) {
	ctx := context.Background()

	t.Run("happy path", func(t *testing.T) {
		// given
		ramcacher, err := cacher.NewRAMCacher()
		require.NoError(t, err)
		uri := "https://example.com/data"
		content := []byte("hello world")

		// when
		err = ramcacher.Store(ctx, uri, content)
		require.NoError(t, err)
		retrieved, err := ramcacher.Retrieve(ctx, uri)

		// then
		require.NoError(t, err)
		assert.Equal(t, content, retrieved)
	})

	t.Run("retrieve uri not found", func(t *testing.T) {
		// given
		ramcacher, _ := cacher.NewRAMCacher()
		uri := "non-existent"

		// when
		_, err := ramcacher.Retrieve(ctx, uri)

		// then
		require.Error(t, err)
		assert.True(t, errors.Is(err, cacher.ErrRAMCacher))
		assert.True(t, errors.Is(err, cacher.ErrURINotFound))
		assert.Contains(t, err.Error(), uri)
	})

	t.Run("close clears the vault", func(t *testing.T) {
		// given
		ramcacher, _ := cacher.NewRAMCacher()
		uri := "key"
		data := []byte("value")
		err := ramcacher.Store(ctx, uri, data)
		require.NoError(t, err)

		got, err := ramcacher.Retrieve(ctx, uri)
		require.NoError(t, err)
		assert.Equal(t, data, got)

		// when
		err = ramcacher.Close()
		require.NoError(t, err)
		_, retrieveErr := ramcacher.Retrieve(ctx, uri)

		// then
		require.Error(t, retrieveErr)
		assert.True(t, errors.Is(retrieveErr, cacher.ErrURINotFound))
	})

	t.Run("initialize with existing vault", func(t *testing.T) {
		// given
		vault := map[string][]byte{
			"pre-existing": []byte("initial"),
		}

		// when
		ramcacher, err := cacher.NewRAMCacherWithVault(vault)
		require.NoError(t, err)
		data, retrieveErr := ramcacher.Retrieve(ctx, "pre-existing")

		// then
		require.NoError(t, retrieveErr)
		assert.Equal(t, []byte("initial"), data)
	})
}
