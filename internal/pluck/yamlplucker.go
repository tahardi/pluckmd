package pluck

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	YAMLIndentSize = 2
)

var (
	ErrYAMLPlucker = errors.New("YAML plucker")
)

type YAMLPlucker struct{}

func NewYAMLPlucker() (*YAMLPlucker, error) {
	return &YAMLPlucker{}, nil
}

func (y *YAMLPlucker) Pluck(
	_ context.Context,
	code string,
	name string,
	kind Kind,
) (string, error) {
	switch kind {
	case File:
		return code, nil
	case Node:
		break
	case Func, Type:
		return "", fmt.Errorf("%w: func and type kind not supported", ErrYAMLPlucker)
	default:
		return "", fmt.Errorf("%w: unrecognized kind: %v", ErrYAMLPlucker, kind)
	}

	node := yaml.Node{}
	err := yaml.Unmarshal([]byte(code), &node)
	if err != nil {
		return "", fmt.Errorf("%w: unmarshaling: %w", ErrYAMLPlucker, err)
	}

	targetNode, targetName, err := FindTargetNode(&node, name)
	if err != nil {
		return "", err
	}

	switch targetNode.Kind {
	case yaml.AliasNode:
		return "", fmt.Errorf("%w: alias nodes are not supported", ErrYAMLPlucker)
	case yaml.DocumentNode:
		return "", fmt.Errorf("%w: document nodes are not supported", ErrYAMLPlucker)
	case yaml.MappingNode:
		return HandleMappingNode(targetNode, targetName)
	case yaml.ScalarNode:
		return HandleScalarNode(targetNode, targetName)
	case yaml.SequenceNode:
		return "", fmt.Errorf("%w: sequence nodes are not supported", ErrYAMLPlucker)
	default:
		return "", fmt.Errorf("%w: unsupported node kind: %v", ErrYAMLPlucker, targetNode.Kind)
	}
}

func FindTargetNode(node *yaml.Node, path string) (*yaml.Node, string, error) {
	switch {
	case node.Kind != yaml.DocumentNode:
		return nil, "", fmt.Errorf("%w: expected document node, got %v", ErrYAMLPlucker, node.Kind)
	case len(node.Content) == 0:
		return nil, "", fmt.Errorf("%w: document node has no content", ErrYAMLPlucker)
	}

	lastKey := ""
	keys := strings.Split(path, ".")
	current := node.Content[0]
	for _, key := range keys {
		if current.Kind != yaml.MappingNode {
			return nil, "", fmt.Errorf(
				"%w: expected mapping node at key '%s', got %v",
				ErrYAMLPlucker, key, current.Kind,
			)
		}

		// YAML mappings store key-value pairs as consecutive elements in Content
		// Index 0, 2, 4... are keys; Index 1, 3, 5... are values
		found := false
		for i := 0; i < len(current.Content)-1; i += 2 {
			if current.Content[i].Value == key {
				current = current.Content[i+1]
				lastKey = key
				found = true
				break
			}
		}

		if !found {
			return nil, "", fmt.Errorf(
				"%w: key '%s' not found in YAML",
				ErrYAMLPlucker, key,
			)
		}
	}
	return current, lastKey, nil
}

func HandleScalarNode(node *yaml.Node, name string) (string, error) {
	if node.Kind != yaml.ScalarNode {
		return "", fmt.Errorf(
			"%w: expected scalar node at key '%s', got %v",
			ErrYAMLPlucker, name, node.Kind,
		)
	}

	body, err := YAMLBodyWithIndent(node)
	if err != nil {
		return "", fmt.Errorf(
			"%w: building YAML body for key '%s': %w",
			ErrYAMLPlucker, name, err,
		)
	}
	return fmt.Sprintf("%s: %s", name, body), nil
}

func HandleMappingNode(node *yaml.Node, name string) (string, error) {
	if node.Kind != yaml.MappingNode {
		return "", fmt.Errorf(
			"%w: expected mapping node at key '%s', got %v",
			ErrYAMLPlucker, name, node.Kind,
		)
	}

	body, err := YAMLBodyWithIndent(node)
	if err != nil {
		return "", fmt.Errorf(
			"%w: building YAML body for key '%s': %w",
			ErrYAMLPlucker, name, err,
		)
	}

	lines := strings.Split(body, "\n")
	var indented strings.Builder
	for _, line := range lines {
		if line != "" {
			indented.WriteString("  " + line + "\n")
		}
	}
	return fmt.Sprintf("%s:\n%s", name, indented.String()), nil
}

func YAMLBodyWithIndent(node *yaml.Node) (string, error) {
	var body bytes.Buffer
	encoder := yaml.NewEncoder(&body)
	defer encoder.Close()

	encoder.SetIndent(YAMLIndentSize)
	err := encoder.Encode(node)
	if err != nil {
		return "", err
	}
	return body.String(), nil
}
