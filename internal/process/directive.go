package process

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/tahardi/pluckmd/internal/pluck"
)

const (
	LangIndex   = 1
	KindIndex   = 2
	NameIndex   = 3
	SourceIndex = 4
	StartIndex  = 5
	EndIndex    = 6
	NumFields   = 7
)

var (
	ErrDirective = errors.New("pluck")

	// Regex subcomponents for building PluckRegex
	commentStart = `<!--`
	commentEnd   = `-->`
	comma        = `,`
	optionalWs   = `\s*`
	pluckName    = `pluck`
	number       = `(-?\d+)`
	quotedString = `"([^"]+)"`

	// PluckRegex The pluck directive we look for in md files is of the form:
	// <!-- pluck("lang", "kind", "name", "source", start, end) -->
	//
	//	lang = "go", "yaml", etc.
	//	kind = "file", "function", "type", etc.
	//	name = name of the type/function/node
	//	source = relative path for local file or remote git URL
	//  start = integer representing starting line of code body
	//  end = integer representing ending line of code body
	// This regex will match with the directive defined above
	PluckRegex = regexp.MustCompile(
		commentStart + optionalWs + pluckName + `\(` +
			optionalWs + quotedString + optionalWs + comma +
			optionalWs + quotedString + optionalWs + comma +
			optionalWs + quotedString + optionalWs + comma +
			optionalWs + quotedString + optionalWs + comma +
			optionalWs + number + optionalWs + comma +
			optionalWs + number + optionalWs +
			`\)` + optionalWs + commentEnd,
	)
)

func ContainsPluckDirective(line string) bool {
	fields := PluckRegex.FindStringSubmatch(line)
	return len(fields) == NumFields
}

type Directive struct {
	lang   pluck.Lang
	kind   pluck.Kind
	name   string
	source string
	start  int
	end    int
}

func NewDirective(line string) (*Directive, error) {
	fields := PluckRegex.FindStringSubmatch(line)
	if len(fields) != NumFields {
		return nil, fmt.Errorf("%w: directive incorrect num fields: %s", ErrDirective, line)
	}

	start, err := strconv.Atoi(fields[StartIndex])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid start index: %s", ErrDirective, line)
	}
	end, err := strconv.Atoi(fields[EndIndex])
	if err != nil {
		return nil, fmt.Errorf("%w: invalid end index: %s", ErrDirective, line)
	}

	lang := pluck.Lang(fields[LangIndex])
	if !lang.Valid() {
		return nil, fmt.Errorf("%w: invalid lang: %s", ErrDirective, line)
	}

	kind := pluck.Kind(fields[KindIndex])
	if !kind.Valid() {
		return nil, fmt.Errorf("%w: invalid kind: %s", ErrDirective, line)
	}

	return &Directive{
		lang:   lang,
		kind:   kind,
		name:   fields[NameIndex],
		source: fields[SourceIndex],
		start:  start,
		end:    end,
	}, nil
}

func (d *Directive) Lang() pluck.Lang {
	return d.lang
}

func (d *Directive) Kind() pluck.Kind {
	return d.kind
}

func (d *Directive) Name() string {
	return d.name
}

func (d *Directive) Start() int {
	return d.start
}

func (d *Directive) End() int {
	return d.end
}

func (d *Directive) CodeSnippetURI() string {
	return d.source + "." + string(d.kind) + "." + d.name
}

func (d *Directive) SourceCodeURI() string {
	return d.source
}
