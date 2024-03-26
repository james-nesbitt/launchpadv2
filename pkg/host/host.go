package host

import (
	"context"
	"io"
)

type Host interface {
	Exec(ctx context.Context, cmd string, inr io.Reader) (string, string, error)
	HasRole(string) bool
}
