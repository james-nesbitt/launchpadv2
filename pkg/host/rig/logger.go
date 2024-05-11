package righost

import (
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/host"
)

// LogCmdLogger write slogs.
type LogCmdLogger struct {
	host  host.Host
	level string
}

func (rlcl LogCmdLogger) Write(p []byte) (n int, err error) {
	msg := fmt.Sprintf("%s: %s", rlcl.host.Id(), string(p))

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
