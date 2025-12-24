package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/tahardi/pluckmd/internal/cache"
	"github.com/tahardi/pluckmd/internal/fetch"
	"github.com/tahardi/pluckmd/internal/pluck"
)

const (
	CodeBlockStartLine = "```go\n"
	CodeBlockStopLine  = "```\n"
)

var (
	ErrProcessor             = errors.New("processor")
	ErrCodeBlockStopNotFound = errors.New("finding end of code block")
)

type Processor struct {
	cacher  cache.Cacher
	fetcher fetch.Fetcher
	plucker pluck.Plucker
}

func NewProcessor(
	cacher cache.Cacher,
	fetcher fetch.Fetcher,
	plucker pluck.Plucker,
) *Processor {
	return &Processor{cacher: cacher, fetcher: fetcher, plucker: plucker}
}

func (p *Processor) ProcessMarkdown(
	ctx context.Context,
	md []byte,
) ([]byte, error) {
	var processed bytes.Buffer
	lines := strings.Split(string(md), "\n")
	for i := 0; i < len(lines); i++ {
		// If it's the last element and it's empty, it's just the trailing newline from the file
		if i == len(lines)-1 && lines[i] == "" {
			break
		}

		processed.WriteString(lines[i] + "\n")
		if !ContainsPluckDirective(lines[i]) {
			continue
		}

		directiveLine := lines[i]
		directive, err := NewDirective(directiveLine)
		if err != nil {
			return nil, fmt.Errorf("%w: creating directive: %w", ErrProcessor, err)
		}

		snippet, err := p.GetCodeSnippet(ctx, directive)
		if err != nil {
			return nil, err
		}

		code, err := snippet.Partial(directive.start, directive.end)
		if err != nil {
			return nil, fmt.Errorf("%w: getting partial code: %w", ErrProcessor, err)
		}

		err = WriteCodeBlock(&processed, directiveLine, code)
		if err != nil {
			return nil, fmt.Errorf("%w: writing code block: %w", ErrProcessor, err)
		}

		end, err := FindCodeBlockEnd(lines, i)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: %w: snippet uri: %s",
				ErrProcessor,
				err,
				directive.CodeSnippetURI(),
			)
		}
		i = end
	}
	return processed.Bytes(), nil
}

func (p *Processor) GetCodeSnippet(
	ctx context.Context,
	directive *Directive,
) (*pluck.Snippet, error) {
	snippetBytes, err := p.cacher.Retrieve(ctx, directive.CodeSnippetURI())
	switch {
	case err == nil:
		return pluck.NewSnippet(directive.Name(), string(snippetBytes))
	case errors.Is(err, cache.ErrURINotFound):
		break
	default:
		return nil, fmt.Errorf(
			"%w: retrieving snippet bytes: %w",
			ErrProcessor,
			err,
		)
	}

	sourceCode, err := p.GetSourceCode(ctx, directive)
	if err != nil {
		return nil, err
	}

	snippetString, err := p.plucker.Pluck(
		string(sourceCode),
		directive.Name(),
		directive.Kind(),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: plucking snippet: %w", ErrProcessor, err)
	}

	err = p.cacher.Store(ctx, directive.CodeSnippetURI(), []byte(snippetString))
	if err != nil {
		return nil, fmt.Errorf(
			"%w: storing snippet bytes: %w",
			ErrProcessor,
			err,
		)
	}
	return pluck.NewSnippet(directive.Name(), snippetString)
}

func (p *Processor) GetSourceCode(
	ctx context.Context,
	directive *Directive,
) ([]byte, error) {
	sourceCode, err := p.cacher.Retrieve(ctx, directive.SourceCodeURI())
	switch {
	case err == nil:
		return sourceCode, nil
	case errors.Is(err, cache.ErrURINotFound):
		break
	default:
		return nil, fmt.Errorf(
			"%w: retrieving source bytes: %w",
			ErrProcessor,
			err,
		)
	}

	sourceCode, err = p.fetcher.Fetch(ctx, directive.SourceCodeURI())
	if err != nil {
		return nil, fmt.Errorf(
			"%w: fetching source bytes: %w",
			ErrProcessor,
			err,
		)
	}

	err = p.cacher.Store(ctx, directive.SourceCodeURI(), sourceCode)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: storing source bytes: %w",
			ErrProcessor,
			err,
		)
	}
	return sourceCode, nil
}

func WriteCodeBlock(
	processed *bytes.Buffer,
	directiveLine string,
	code string,
) error {
	// Calculate indentation by trimming whitespace from the left side of the
	// directive line. Technically, this will match spaces, tabs, and newlines.
	// For example:
	//
	// line = "\t\t<!-- pluck(...) -->", indent = "\t\t"
	// line = " \t \t<!-- pluck(...) -->", indent = " \t \t"
	indent := ""
	trimmed := strings.TrimLeftFunc(directiveLine, unicode.IsSpace)
	if len(trimmed) < len(directiveLine) {
		indent = directiveLine[:len(directiveLine)-len(trimmed)]
	}

	if indent == "" {
		processed.WriteString(CodeBlockStartLine)
		processed.WriteString(code)
		processed.WriteString(CodeBlockStopLine)
		return nil
	}

	indentedCode := IndentCode(code, indent)
	processed.WriteString(indent + CodeBlockStartLine)
	processed.WriteString(indentedCode)
	processed.WriteString(indent + CodeBlockStopLine)
	return nil
}

func IndentCode(code string, indentation string) string {
	if indentation == "" {
		return code
	}

	// Remove newline at end of code string. Otherwise, split will produce
	// a lines array where the last element is an empty line
	lines := strings.Split(strings.TrimSuffix(code, "\n"), "\n")
	for i, line := range lines {
		lines[i] = indentation + line
	}
	return strings.Join(lines, "\n") + "\n"
}

func FindCodeBlockEnd(lines []string, i int) (int, error) {
	foundStart := false
	for ; i < len(lines); i++ {
		trimmed := strings.TrimLeftFunc(lines[i], unicode.IsSpace) + "\n"
		if !foundStart {
			if strings.HasPrefix(trimmed, CodeBlockStartLine) {
				foundStart = true
			}
		} else {
			if strings.HasSuffix(trimmed, CodeBlockStopLine) {
				return i, nil
			}
		}
	}
	return len(lines), ErrCodeBlockStopNotFound
}
