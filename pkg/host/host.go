package host

import (
	"context"
	"io"
)

type Host interface {
	Id() string
	HasRole(string) bool

	Exec(ctx context.Context, cmd string, inr io.Reader) (string, string, error)

	IsWindows() bool
}
