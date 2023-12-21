package phase

import "context"

type Phase interface {
	Run(context.Context) error
}
