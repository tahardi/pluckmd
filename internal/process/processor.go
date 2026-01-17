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
	"github.com/tahardi/pluckmd/internal/snip"
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
	cacher   cache.Cacher
	fetchers []fetch.Fetcher
	pluckers map[pluck.Lang]pluck.Plucker
}

func NewProcessor(
	cacher cache.Cacher,
	fetchers []fetch.Fetcher,
	pluckers map[pluck.Lang]pluck.Plucker,
) *Processor {
	return &Processor{cacher: cacher, fetchers: fetchers, pluckers: pluckers}
}

func (p *Processor) ProcessMarkdown(
	ctx context.Context,
	md []byte,
) ([]byte, error) {
	// Split markdown into lines. If the markdown ends with a newline, Split
	// will return an empty string as the last element. This will cause us
	// to write an extra newline to the output. Remove the ending newline
	// (if any) so that we don't write an extra newline.
	lines := strings.Split(strings.TrimSuffix(string(md), "\n"), "\n")

	var processed bytes.Buffer
	for i := 0; i < len(lines); i++ {
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
			return nil, fmt.Errorf("%w: getting snippet: %w", ErrProcessor, err)
		}

		err = WriteCodeBlock(&processed, directiveLine, snippet)
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
) (string, error) {
	fullSnippet, err := p.GetFullCodeSnippet(ctx, directive)
	if err != nil {
		return "", err
	}

	var snipper snip.Snipper
	switch directive.Lang() {
	case pluck.Go:
		snipper, err = snip.NewGoSnipper(directive.Name(), fullSnippet)
		if err != nil {
			return "", fmt.Errorf("%w: creating go snipper: %w", ErrProcessor, err)
		}
	case pluck.YAML:
		snipper, err = snip.NewYAMLSnipper(directive.Name(), fullSnippet)
		if err != nil {
			return "", fmt.Errorf("%w: creating yaml snipper: %w", ErrProcessor, err)
		}
	default:
		return "", fmt.Errorf("%w: unsupported lang: %s", ErrProcessor, directive.Lang())
	}
	return snipper.Snippet(directive.start, directive.end)
}

func (p *Processor) GetFullCodeSnippet(
	ctx context.Context,
	directive *Directive,
) (string, error) {
	snippetBytes, err := p.cacher.Retrieve(ctx, directive.CodeSnippetURI())
	switch {
	case err == nil:
		return string(snippetBytes), nil
	case errors.Is(err, cache.ErrURINotFound):
		break
	default:
		return "", fmt.Errorf(
			"%w: retrieving snippet bytes: %w",
			ErrProcessor,
			err,
		)
	}

	sourceCode, err := p.GetSourceCode(ctx, directive)
	if err != nil {
		return "", err
	}

	plucker, exists := p.pluckers[directive.Lang()]
	if !exists {
		return "", fmt.Errorf("%w: no plucker for lang: %s", ErrProcessor, directive.Lang())
	}

	snippetString, err := plucker.Pluck(
		ctx,
		string(sourceCode),
		directive.Name(),
		directive.Kind(),
	)
	if err != nil {
		return "", fmt.Errorf("%w: plucking snippet: %w", ErrProcessor, err)
	}

	err = p.cacher.Store(ctx, directive.CodeSnippetURI(), []byte(snippetString))
	if err != nil {
		return "", fmt.Errorf(
			"%w: storing snippet bytes: %w",
			ErrProcessor,
			err,
		)
	}
	return snippetString, nil
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

	sourceCode, err = p.Fetch(ctx, directive)
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

func (p *Processor) Fetch(ctx context.Context, directive *Directive) ([]byte, error) {
	errs := []error{}
	for _, fetcher := range p.fetchers {
		sourceCode, err := fetcher.Fetch(ctx, directive.SourceCodeURI())
		if err == nil {
			return sourceCode, nil
		}
		errs = append(errs, fmt.Errorf(
			"%w: fetching source bytes: %w",
			ErrProcessor,
			err),
		)
	}
	return nil, errors.Join(errs...)
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
