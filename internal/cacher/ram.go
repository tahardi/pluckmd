package cacher

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrRAMCacher   = errors.New("ram cacher")
	ErrURINotFound = errors.New("uri not found")
)

type RAMCacher struct {
	vault map[string][]byte
}

func NewRAMCacher() (*RAMCacher, error) {
	return NewRAMCacherWithVault(make(map[string][]byte))
}

func NewRAMCacherWithVault(vault map[string][]byte) (*RAMCacher, error) {
	return &RAMCacher{vault: vault}, nil
}

func (r *RAMCacher) Store(_ context.Context, uri string, data []byte) error {
	r.vault[uri] = data
	return nil
}

func (r *RAMCacher) Retrieve(_ context.Context, uri string) ([]byte, error) {
	data, ok := r.vault[uri]
	if ok {
		return data, nil
	}
	return nil, fmt.Errorf("%w: %w: %s", ErrRAMCacher, ErrURINotFound, uri)
}

func (r *RAMCacher) Close() error {
	r.vault = map[string][]byte{}
	return nil
}
