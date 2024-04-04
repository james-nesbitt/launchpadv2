package righost

import (
	"context"
	"io"
	"log/slog"
	"strings"
)

// NewRigHostHost Host constructor from HostConfig
func NewHost(hc *Config) Host {
	return Host{
		Config: hc,
		state:  &State{},
	}
}

// Host definition for one cluster target
type Host struct {
	*Config
	state *State
}

// HasRole answer boolean if a host has a role
func (h Host) HasRole(role string) bool {
	return h.Config.Role == role
}

// Exec a command
func (h *Host) Exec(ctx context.Context, cmd string, inr io.Reader) (string, string, error) {
	outs := strings.Builder{}
	outw := io.MultiWriter(
		LogCmdLogger{level: "out"},
		&outs,
	)

	errs := strings.Builder{}
	errw := io.MultiWriter(
		LogCmdLogger{level: "err"},
		&errs,
	)

	w, serr := h.StartProcess(ctx, cmd, inr, outw, errw)
	if serr != nil {
		return outs.String(), errs.String(), serr
	}

	if werr := w.Wait(); werr != nil {
		return outs.String(), errs.String(), werr
	}

	return outs.String(), errs.String(), nil
}

// RigHostState state for the host
type State struct {
}

// LogCmdLogger write slogs
type LogCmdLogger struct {
	level string
}

func (rlcl LogCmdLogger) Write(p []byte) (n int, err error) {
	if rlcl.level == "err" {
		slog.Debug(string(p))
	} else {
		slog.Warn(string(p))
	}
	return len(p), nil
}