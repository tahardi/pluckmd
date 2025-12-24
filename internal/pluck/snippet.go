package pluck

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const (
	EmptyStart   = -1
	EmptyEnd     = -1
	FullStart    = 0
	FullEnd      = 0
	EllipsesLine = "\t// ...\n"
	OpeningBrace = "{\n"
	ClosingBrace = "}\n"
)

var (
	ErrSnippet              = errors.New("snippet")
	ErrOpeningBraceNotFound = errors.New("finding opening brace")
)

type Snippet struct {
	name       string
	definition string
	body       string
	bodyLines  []string
	length     int
}

func NewSnippet(name string, snippet string) (*Snippet, error) {
	definition, body, err := ParseSnippet(snippet)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing snippet: %w", ErrSnippet, err)
	}
	return NewSnippetWithDefinitionAndBody(name, definition, body)
}

func NewSnippetWithDefinitionAndBody(
	name string,
	definition string,
	body string,
) (*Snippet, error) {
	return &Snippet{
		name:       name,
		definition: definition,
		body:       body,
		bodyLines:  nil,
		length:     0,
	}, nil
}

func (s *Snippet) Name() string {
	return s.name
}

func (s *Snippet) Definition() string {
	return s.definition
}

func (s *Snippet) Body() string {
	return s.body
}

func (s *Snippet) Full() string {
	return s.definition + OpeningBrace + s.body + ClosingBrace
}

func (s *Snippet) Empty() string {
	return s.definition + OpeningBrace + EllipsesLine + ClosingBrace
}

func (s *Snippet) Partial(start int, end int) (string, error) {
	// Lazy initialization of bodyLines. Sometimes we end up with empty lines
	// at the beginning and ending of the body. If so, remove them.
	if s.bodyLines == nil {
		lines := strings.Split(s.body, "\n")
		if len(lines) > 0 {
			if lines[0] == "" {
				lines = lines[1:]
			}
			if lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
		}
		s.bodyLines = lines
		s.length = len(lines)
	}

	switch {
	case start == EmptyStart && end == EmptyEnd:
		return s.Empty(), nil
	case start == FullStart && end == FullEnd:
		return s.Full(), nil
	case start < 0 || end < 0 || start > end || end > s.length:
		return "", fmt.Errorf(
			"%w: invalid range [start: %d, end: %d)",
			ErrSnippet,
			start,
			end,
		)
	}

	// If we are skipping the beginning of the body, add an Ellipses line to
	// indicate that there is hidden code we are not including.
	snippet := s.Definition() + OpeningBrace
	if start != 0 {
		snippet += EllipsesLine
	}

	for i := start; i < end; i++ {
		snippet += s.bodyLines[i] + "\n"
	}

	// If we are skipping the end of the body, add an Ellipses line to indicate
	// that there is hidden code we are not including.
	if end != s.length {
		snippet += EllipsesLine
	}
	snippet += ClosingBrace
	return snippet, nil
}

func ParseSnippet(snippet string) (string, string, error) {
	fset := token.NewFileSet()

	// Wrap snippet in a package to make it a valid Go file for the parser
	dummyPackage := "package dummy\n"
	dummySource := dummyPackage + snippet
	f, err := parser.ParseFile(fset, "", dummySource, 0)
	if err != nil {
		return "", "", fmt.Errorf("making ast file: %w", err)
	}

	var openingBracePos token.Pos
	var closingBracePos token.Pos

	// Find the first function or type declaration
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Body != nil {
				openingBracePos = x.Body.Lbrace
				closingBracePos = x.Body.Rbrace
				return false
			}
		case *ast.TypeSpec:
			if st, ok := x.Type.(*ast.StructType); ok {
				openingBracePos = st.Fields.Opening
				closingBracePos = st.Fields.Closing
				return false
			}
			if it, ok := x.Type.(*ast.InterfaceType); ok {
				openingBracePos = it.Methods.Opening
				closingBracePos = it.Methods.Closing
				return false
			}
		}
		return true
	})

	if !openingBracePos.IsValid() {
		return "", "", ErrOpeningBraceNotFound
	}

	// Calculate offsets relative to the original code string
	// Subtract 1 because we added "package dummy\n" (14 chars)
	offset := fset.Position(openingBracePos).Offset - len(dummyPackage)
	endOffset := fset.Position(closingBracePos).Offset - len(dummyPackage)

	definition := snippet[:offset]
	body := snippet[offset+1 : endOffset]
	if body[0] == '\n' {
		body = body[1:]
	}
	return definition, body, nil
}
