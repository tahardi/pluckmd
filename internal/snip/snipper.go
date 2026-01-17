package snip

const (
	EmptyStart   = -1
	EmptyEnd     = -1
	FullStart    = 0
	FullEnd      = 0
)

type Snipper interface {
	Snippet(start int, end int) (string, error)
}
