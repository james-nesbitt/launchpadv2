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
	if rlcl.level == "err" {
		slog.Warn(msg)
	} else {
		slog.Info(msg)
	}
	return len(p), nil
}
