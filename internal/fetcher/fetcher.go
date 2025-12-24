package fetch

import "context"

type Fetcher interface {
	Fetch(ctx context.Context, uri string) (data []byte, err error)
}
