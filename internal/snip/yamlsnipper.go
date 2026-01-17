package snip

import (
	"errors"
)

var (
	ErrYAMLSnipper = errors.New("YAMLSnipper error")
)

type YAMLSnipper struct {
	name    string
	snippet string
}

func NewYAMLSnipper(name string, snippet string) (*YAMLSnipper, error) {
	return &YAMLSnipper{
		name:    name,
		snippet: snippet,
	}, nil
}

func (y *YAMLSnipper) Snippet(_ int, _ int) (string, error) {
	return y.snippet, nil
}
