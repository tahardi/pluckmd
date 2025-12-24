package plucker_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tahardi/pluckmd/internal/plucker"
)

const (
	blockyPlucker             = "BlockyPlucker"
	blockyPluckerSnippet      = `type BlockyPlucker struct{}`
	blockyPluckerPluck        = "BlockyPlucker.Pluck"
	blockyPluckerPluckSnippet = `func (b *BlockyPlucker) Pluck(
	code string,
	name string,
	kind Kind,
) (string, error) {
	if !kind.Valid() {
		return "", fmt.Errorf("%w: invalid kind '%s'", ErrBlockyPlucker, kind)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	pick := fmt.Sprintf("%s=%s:%s", PickArg, kind, name)

	cmd := exec.Command(PluckCmd, pick)
	cmd.Stdin = strings.NewReader(code)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

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
}`
)

//go:embed blocky.go
var blockyGo string

func TestBlockyPlucker_Pluck(t *testing.T) {
	t.Run("happy path - BlockyPlucker (type)", func(t *testing.T) {
		// given
		code := blockyGo
		name := blockyPlucker
		kind := plucker.Type
		want := blockyPluckerSnippet + "\n"
		blockyplucker, err := plucker.NewBlockyPlucker()
		require.NoError(t, err)

		// when
		got, err := blockyplucker.Pluck(code, name, kind)

		// then
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("happy path - BlockyPlucker.Pluck (func)", func(t *testing.T) {
		code := blockyGo
		name := blockyPluckerPluck
		kind := plucker.Func
		want := blockyPluckerPluckSnippet + "\n"
		blockyplucker, err := plucker.NewBlockyPlucker()
		require.NoError(t, err)

		// when
		got, err := blockyplucker.Pluck(code, name, kind)

		// then
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("error - invalid kind", func(t *testing.T) {
		code := blockyGo
		name := blockyPluckerPluck
		kind := plucker.Kind("invalid")
		blockyplucker, err := plucker.NewBlockyPlucker()
		require.NoError(t, err)

		// when
		_, err = blockyplucker.Pluck(code, name, kind)

		// then
		require.Error(t, err)
	})

	t.Run("error - type/func not in code", func(t *testing.T) {
		code := blockyGo
		name := "funcDoesNotExist"
		kind := plucker.Func
		blockyplucker, err := plucker.NewBlockyPlucker()
		require.NoError(t, err)

		// when
		_, err = blockyplucker.Pluck(code, name, kind)

		// then
		require.Error(t, err)
	})
}
