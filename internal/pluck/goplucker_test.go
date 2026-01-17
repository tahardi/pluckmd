package pluck_test

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tahardi/pluckmd/internal/pluck"
)

const (
	goPlucker             = "GoPlucker"
	goPluckerSnippet      = `type GoPlucker struct{}`
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

//go:embed goplucker.go
var goPluckerSource string

func TestGoPlucker_Pluck(t *testing.T) {
	t.Run("happy path - GoPlucker (type)", func(t *testing.T) {
		// given
		ctx := context.Background()
		code := goPluckerSource
		name := goPlucker
		kind := pluck.Type
		want := goPluckerSnippet + "\n"
		plucker, err := pluck.NewGoPlucker()
		require.NoError(t, err)

		// when
		got, err := plucker.Pluck(ctx, code, name, kind)

		// then
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("happy path - GoPlucker.Pluck (func)", func(t *testing.T) {
		// given
		ctx := context.Background()
		code := goPluckerSource
		name := goPluckerPluck
		kind := pluck.Func
		want := goPluckerPluckSnippet + "\n"
		plucker, err := pluck.NewGoPlucker()
		require.NoError(t, err)

		// when
		got, err := plucker.Pluck(ctx, code, name, kind)

		// then
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("error - unsupported kind", func(t *testing.T) {
		// given
		ctx := context.Background()
		code := goPluckerSource
		name := goPluckerPluck
		kind := pluck.Node
		plucker, err := pluck.NewGoPlucker()
		require.NoError(t, err)

		// when
		_, err = plucker.Pluck(ctx, code, name, kind)

		// then
		require.Error(t, err)
	})

	t.Run("error - type/func not in code", func(t *testing.T) {
		// given
		ctx := context.Background()
		code := goPluckerSource
		name := "funcDoesNotExist"
		kind := pluck.Func
		plucker, err := pluck.NewGoPlucker()
		require.NoError(t, err)

		// when
		_, err = plucker.Pluck(ctx, code, name, kind)

		// then
		require.Error(t, err)
	})

	t.Run("error - context canceled", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		code := goPluckerSource
		name := goPluckerPluck
		kind := pluck.Func
		plucker, err := pluck.NewGoPlucker()
		require.NoError(t, err)

		// when
		_, err = plucker.Pluck(ctx, code, name, kind)

		// then
		require.Error(t, err)
	})
}
