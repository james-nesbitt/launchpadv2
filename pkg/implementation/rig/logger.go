package rig

import (
	"fmt"
	"log/slog"
)

// LogCmdLogger write slogs.
type LogCmdLogger struct {
	id    string
	level string
}

func (rlcl LogCmdLogger) Write(p []byte) (n int, err error) {
	msg := fmt.Sprintf("%s: %s", rlcl.id, string(p))

	switch rlcl.level {

	case "info":
		slog.Info(msg)
	case "debug":
		slog.Debug(msg)
	case "warn":
		slog.Warn(msg)
	case "error":
		slog.Error(msg)
	default:
		slog.Debug(msg)
	}
	return len(p), nil
}
