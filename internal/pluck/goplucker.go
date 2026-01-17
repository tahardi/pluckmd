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
	GoPluckCmd       = "pluck"
	GoPluckCLISource = "github.com/blocky/pluck/cmd/pluck@v0.1.1"
	PickArg          = "--pick"
)

var (
	ErrGoPlucker        = errors.New("go plucker")
	ErrPluckCmdNotFound = fmt.Errorf("pluck command '%s' not found", GoPluckCmd)
)

type GoPlucker struct{}

func NewGoPlucker() (*GoPlucker, error) {
	_, err := exec.LookPath(GoPluckCmd)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %w: install '%s' via 'go install %s'",
			ErrGoPlucker,
			ErrPluckCmdNotFound,
			GoPluckCmd,
			GoPluckCLISource,
		)
	}
	return &GoPlucker{}, nil
}

func (g *GoPlucker) Pluck(
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
}
