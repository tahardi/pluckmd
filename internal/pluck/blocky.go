package pluck

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	PluckCmd       = "pluck"
	PluckCLISource = "github.com/blocky/pluck/cmd/pluck@v0.1.1"
	PickArg        = "--pick"
)

var (
	ErrBlockyPlucker    = errors.New("blocky plucker")
	ErrPluckCmdNotFound = fmt.Errorf("pluck command '%s' not found", PluckCmd)
)

type BlockyPlucker struct{}

func NewBlockyPlucker() (*BlockyPlucker, error) {
	_, err := exec.LookPath(PluckCmd)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %w: install '%s' via 'go install %s'",
			ErrBlockyPlucker,
			ErrPluckCmdNotFound,
			PluckCmd,
			PluckCLISource,
		)
	}
	return &BlockyPlucker{}, nil
}

func (b *BlockyPlucker) Pluck(
	ctx context.Context,
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

	cmd := exec.CommandContext(ctx, PluckCmd, pick)
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
}
