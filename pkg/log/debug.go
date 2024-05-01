package log

import "context"

// Debugable can debug output a description of itself.
type Debuggable interface {
	Debug(context.Context) interface{}
}
