package types

import "context"

type Cleaner interface {
	Clean(context.Context, []byte) ([]byte, error)
}
