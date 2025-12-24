package pluck

import "context"

type Plucker interface {
	Pluck(ctx context.Context, code string, name string, kind Kind) (snippet string, err error)
}
