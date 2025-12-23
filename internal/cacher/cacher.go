package cacher

import "context"

type Cacher interface {
	Close() error
	Store(ctx context.Context, uri string, data []byte) (err error)
	Retrieve(ctx context.Context, uri string) (data []byte, err error)
}