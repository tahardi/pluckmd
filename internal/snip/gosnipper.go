package snip

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const (
	EllipsesLine = "\t// ...\n"
	OpeningBrace = "{\n"
	ClosingBrace = "}\n"
)

var (
	ErrGoSnipper            = errors.New("go snipper")
	ErrOpeningBraceNotFound = errors.New("finding opening brace")
)

type GoSnipper struct {
	name       string
	definition string
	body       string
	bodyLines  []string
	length     int
}

func NewGoSnipper(name string, snippet string) (*GoSnipper, error) {
	definition, body, err := ParseGoSnippet(snippet)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing snippet: %w", ErrGoSnipper, err)
	}
	return NewGoSnipperWithDefinitionAndBody(name, definition, body)
}

func NewGoSnipperWithDefinitionAndBody(
	name string,
	definition string,
	body string,
) (*GoSnipper, error) {
	return &GoSnipper{
		name:       name,
		definition: definition,
		body:       body,
		bodyLines:  nil,
		length:     0,
	}, nil
}

func (g *GoSnipper) Name() string {
	return g.name
}

func (g *GoSnipper) Definition() string {
	return g.definition
}

func (g *GoSnipper) Body() string {
	return g.body
}

func (g *GoSnipper) Full() string {
	return g.definition + OpeningBrace + g.body + ClosingBrace
}

func (g *GoSnipper) Empty() string {
	return g.definition + OpeningBrace + EllipsesLine + ClosingBrace
}

func (g *GoSnipper) Snippet(start int, end int) (string, error) {
	// Lazy initialization of bodyLines. Sometimes we end up with empty lines
	// at the beginning and ending of the body. If so, remove them.
	if g.bodyLines == nil {
		lines := strings.Split(g.body, "\n")
		if len(lines) > 0 {
			if lines[0] == "" {
				lines = lines[1:]
			}
			if lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
		}
		g.bodyLines = lines
		g.length = len(lines)
	}

	switch {
	case start == EmptyStart && end == EmptyEnd:
		return g.Empty(), nil
	case start == FullStart && end == FullEnd:
		return g.Full(), nil
	case start < 0 || end < 0 || start > end || end > g.length:
		return "", fmt.Errorf(
			"%w: invalid range [start: %d, end: %d)",
			ErrGoSnipper,
			start,
			end,
		)
	}

	// If we are skipping the beginning of the body, add an Ellipses line to
	// indicate that there is hidden code we are not including.
	var snippet strings.Builder
	snippet.WriteString(g.Definition())
	snippet.WriteString(OpeningBrace)
	if start != 0 {
		snippet.WriteString(EllipsesLine)
	}

	for i := start; i < end; i++ {
		snippet.WriteString(g.bodyLines[i])
		snippet.WriteString("\n")
	}

	// If we are skipping the end of the body, add an Ellipses line to indicate
	// that there is hidden code we are not including.
	if end != g.length {
		snippet.WriteString(EllipsesLine)
	}
	snippet.WriteString(ClosingBrace)
	return snippet.String(), nil
}

func ParseGoSnippet(snippet string) (string, string, error) {
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
