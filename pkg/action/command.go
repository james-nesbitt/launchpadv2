package action

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Mirantis/launchpad/pkg/dependency"
)

const (
	// Apply command means ensure that resources are in place.
	CommandKeyApply = "apply"
	// Reset command means ensure that resources are removed.
	CommandKeyReset = "reset"
	// Discover command means discover what resources are in place (state).
	CommandKeyDiscover = "discover"
)

var (
	ErrCommandBuild = errors.New("failure building command")
)

// CommandHandler something that wants to participate in building a command.
type CommandHandler interface {
	// Build the command Phases and Events
	CommandBuild(context.Context, *Command) error
}

// NewEmptyCommand create a new Commad that can now be built up.
func NewEmptyCommand(key string) *Command {
	return &Command{
		Key:          key,
		Dependencies: dependency.NewDependencies(),
		Events:       dependency.Events{},
		Phases:       Phases{},
	}
}

// BuildCommand from a number of handler.
//
// @NOTE this looks like it is not used (outside of testing).
func BuildCommand(ctx context.Context, cmd *Command, handlers []CommandHandler) error {
	cberrs := []error{}
	for _, ch := range handlers {
		if ch == nil {
			continue
		}

		if cberr := ch.CommandBuild(ctx, cmd); cberr != nil {
			cberrs = append(cberrs, cberr)
		}
	}
	if len(cberrs) > 0 {
		return fmt.Errorf("%s; command build handle failure: %s", ErrCommandBuild, errors.Join(cberrs...).Error())
	}

	return nil
}

// Command an executable component, typically built by Components.
type Command struct {
	Key string

	Dependencies dependency.Dependencies

	Events dependency.Events
	Phases Phases
}

// Validate the command.
func (cmd *Command) Validate(ctx context.Context) error {
	errs := []error{}

	if cmd.Key == "" {
		errs = append(errs, fmt.Errorf("validate: command in invalid because it has no key"))
	}

	if len(cmd.Phases) == 0 {
		errs = append(errs, fmt.Errorf("validate: command in invalid as it has no phases"))
	} else if pso, oerr := cmd.Phases.Order(ctx); oerr != nil {
		for key, p := range cmd.Phases {
			slog.InfoContext(ctx, fmt.Sprintf("Phase %s", key), slog.Any("phase", p))
		}
		slog.ErrorContext(ctx, "phase order fail")
		errs = append(errs, fmt.Errorf("validate: could not order phases: %w", oerr))
	} else {
		for _, p := range pso {
			pv, ok := p.(Validator)
			if !ok {
				continue
			}

			if err := pv.Validate(ctx); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("command validation failed: %s", errors.Join(errs...).Error())
	}

	return nil
}

// Run execute the command.
//
// @NOTE you should runn all the builder before doing this.
func (cmd *Command) Run(ctx context.Context) error {
	po, oerr := cmd.Phases.Order(ctx)
	if oerr != nil {
		return fmt.Errorf("phase order error: %s", oerr.Error())
	}

	slog.InfoContext(ctx, fmt.Sprintf("COMMAND: %s", cmd.Key), slog.Any("command", cmd))
	for _, p := range po {
		slog.InfoContext(ctx, fmt.Sprintf("PHASE: %s", p.Id()), slog.Any("phase", p))
		if err := p.Run(ctx); err != nil {
			return fmt.Errorf("command run failed on phase [%s] : %s", p.Id(), err.Error())
		}
	}

	return nil
}
