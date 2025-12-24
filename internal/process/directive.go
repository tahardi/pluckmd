package process

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/tahardi/pluckmd/internal/pluck"
)

const (
	KindIndex   = 1
	NameIndex   = 2
	SourceIndex = 3
	StartIndex  = 4
	EndIndex    = 5
	NumFields   = 6
)

var (
	ErrDirective = errors.New("pluck")

	// PluckRegex The pluck directive we look for in md files is of the form:
	// <!-- pluck("kind", "name", "source", start, end) -->
	//
	//	kind = "type" or "function"
	//	name = name of the type or function
	//	source = relative path for local file or remote git URL
	//  start = integer representing starting line of code body
	//  end = integer representing ending line of code body
	// This regex will match with the directive defined above
	PluckRegex = regexp.MustCompile(`<!--\s*pluck\(\s*"([^"]+)"\s*,\s*"([^"]+)"\s*,\s*"([^"]+)"\s*,\s*(-?\d+)\s*,\s*(-?\d+)\s*\)\s*-->`)
)

func ContainsPluckDirective(line string) bool {
	fields := PluckRegex.FindStringSubmatch(line)
	return len(fields) == NumFields
}

type Directive struct {
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

	kind := pluck.Kind(fields[KindIndex])
	if !kind.Valid() {
		return nil, fmt.Errorf("%w: invalid kind: %s", ErrDirective, line)
	}

	return &Directive{
		kind:   kind,
		name:   fields[NameIndex],
		source: fields[SourceIndex],
		start:  start,
		end:    end,
	}, nil
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
