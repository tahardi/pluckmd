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
		want := goPluckerPluckSnippet + "\n"
		snippet, err := pluck.NewSnippet(goPluckerPluck, want)
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
		snippet, err := pluck.NewSnippet(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := pluck.EmptyStart, pluck.EmptyEnd
		want := `func (g *GoPlucker) Pluck(
	ctx context.Context,
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
		snippet, err := pluck.NewSnippet(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := 0, 8
		want := `func (g *GoPlucker) Pluck(
	ctx context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	switch kind {
	case File:
		return code, nil
	case Func, Type:
		break
	default:
		return "", fmt.Errorf("%w: unsupported kind: %v", ErrGoPlucker, kind)
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
		snippet, err := pluck.NewSnippet(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := 9, 17
		want := `func (g *GoPlucker) Pluck(
	ctx context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	// ...
	var out bytes.Buffer
	var stderr bytes.Buffer
	pick := fmt.Sprintf("%s=%s:%s", PickArg, kind, name)

	cmd := exec.CommandContext(ctx, GoPluckCmd, pick)
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
		snippet, err := pluck.NewSnippet(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := 18, 28
		want := `func (g *GoPlucker) Pluck(
	ctx context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	// ...
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf(
			"%w: running %s: %s",
			ErrGoPlucker,
			GoPluckCmd,
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
		snippet, err := pluck.NewSnippet(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := pluck.FullStart, pluck.FullEnd
		want := goPluckerPluckSnippet + "\n"

		// when
		got, err := snippet.Partial(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
