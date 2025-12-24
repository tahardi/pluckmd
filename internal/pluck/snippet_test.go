package pluck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/pluckmd/internal/pluck"
)

func TestSnippet_Full(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		want := blockyPluckerPluckSnippet + "\n"
		snippet, err := pluck.NewSnippet(blockyPluckerPluck, want)
		require.NoError(t, err)

		// when
		got := snippet.Full()

		// then
		require.Equal(t, want, got)
	})
}

func TestSnippet_Partial(t *testing.T) {
	t.Run("happy path - no body", func(t *testing.T) {
		// given
		snippet, err := pluck.NewSnippet(blockyPluckerPluck, blockyPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := pluck.EmptyStart, pluck.EmptyEnd
		want := `func (b *BlockyPlucker) Pluck(
	code string,
	name string,
	kind Kind,
) (string, error) {
	// ...
}
`
		// when
		got, err := snippet.Partial(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - start of body", func(t *testing.T) {
		// given
		snippet, err := pluck.NewSnippet(blockyPluckerPluck, blockyPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := 0, 3
		want := `func (b *BlockyPlucker) Pluck(
	code string,
	name string,
	kind Kind,
) (string, error) {
	if !kind.Valid() {
		return "", fmt.Errorf("%w: invalid kind '%s'", ErrBlockyPlucker, kind)
	}
	// ...
}
`
		// when
		got, err := snippet.Partial(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - middle of body", func(t *testing.T) {
		// given
		snippet, err := pluck.NewSnippet(blockyPluckerPluck, blockyPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := 4, 12
		want := `func (b *BlockyPlucker) Pluck(
	code string,
	name string,
	kind Kind,
) (string, error) {
	// ...
	var out bytes.Buffer
	var stderr bytes.Buffer
	pick := fmt.Sprintf("%s=%s:%s", PickArg, kind, name)

	cmd := exec.Command(PluckCmd, pick)
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	// ...
}
`
		// when
		got, err := snippet.Partial(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - end of body", func(t *testing.T) {
		// given
		snippet, err := pluck.NewSnippet(blockyPluckerPluck, blockyPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := 13, 23
		want := `func (b *BlockyPlucker) Pluck(
	code string,
	name string,
	kind Kind,
) (string, error) {
	// ...
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(
			"%w: running %s: %s",
			ErrBlockyPlucker,
			PluckCmd,
			stderr.String(),
		)
	}
	return out.String(), nil
}
`
		// when
		got, err := snippet.Partial(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - full body", func(t *testing.T) {
		// given
		snippet, err := pluck.NewSnippet(blockyPluckerPluck, blockyPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := pluck.FullStart, pluck.FullEnd
		want := blockyPluckerPluckSnippet + "\n"

		// when
		got, err := snippet.Partial(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
