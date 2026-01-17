package snip_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahardi/pluckmd/internal/snip"
)

const (
	goPluckerPluck        = "GoPlucker.Pluck"
	goPluckerPluckSnippet = `func (g *GoPlucker) Pluck(
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
	case Node:
		return "", fmt.Errorf("%w: node kind not supported", ErrGoPlucker)
	default:
		return "", fmt.Errorf("%w: unrecognized kind: %v", ErrGoPlucker, kind)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	pick := fmt.Sprintf("%s=%s:%s", PickArg, kind, name)

	cmd := exec.CommandContext(ctx, GoPluckCmd, pick)
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

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
}`
)

func TestGoSnippet_Full(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		want := goPluckerPluckSnippet + "\n"
		snippet, err := snip.NewGoSnipper(goPluckerPluck, want)
		require.NoError(t, err)

		// when
		got := snippet.Full()

		// then
		require.Equal(t, want, got)
	})
}

func TestGoSnippet_Partial(t *testing.T) {
	t.Run("happy path - no body", func(t *testing.T) {
		// given
		snippet, err := snip.NewGoSnipper(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := snip.EmptyStart, snip.EmptyEnd
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
		got, err := snippet.Snippet(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - start of body", func(t *testing.T) {
		// given
		snippet, err := snip.NewGoSnipper(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := 0, 10
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
	case Node:
		return "", fmt.Errorf("%w: node kind not supported", ErrGoPlucker)
	default:
		return "", fmt.Errorf("%w: unrecognized kind: %v", ErrGoPlucker, kind)
	}
	// ...
}
`
		// when
		got, err := snippet.Snippet(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - middle of body", func(t *testing.T) {
		// given
		snippet, err := snip.NewGoSnipper(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := 11, 19
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
		got, err := snippet.Snippet(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - end of body", func(t *testing.T) {
		// given
		snippet, err := snip.NewGoSnipper(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := 20, 30
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
		got, err := snippet.Snippet(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("happy path - full body", func(t *testing.T) {
		// given
		snippet, err := snip.NewGoSnipper(goPluckerPluck, goPluckerPluckSnippet)
		require.NoError(t, err)

		start, end := snip.FullStart, snip.FullEnd
		want := goPluckerPluckSnippet + "\n"

		// when
		got, err := snippet.Snippet(start, end)

		// then
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
