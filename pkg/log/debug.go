/*
Package log an attempt to start an instrospective system for the components.
*/
package log

import "context"

// Debuggable can debug output a description of itself.
type Debuggable interface {
	Debug(context.Context) any
}
