package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mirantis/launchpad/pkg/config"
	"github.com/Mirantis/launchpad/pkg/project"
)

var (
	debug   bool
	cfgFile string
	proj    project.Project = project.New()
)

func Bootstrap(root *cobra.Command) project.Project {
	root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	root.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "increase logging verbosity")

	// This should early pre-populate the above flags, which we will use to build more commands. This will
	// be executed again with rootCmd.Execute()
	root.ParseFlags(os.Args[1:])

	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	if cfgFile == "" {
		slog.Warn("No project config specified. Using an empty project")
		proj = project.New()
	} else if err := boostrapBuildProject(root.Context()); err != nil {
		slog.Error("failed to build project", slog.Any("error", err))
		os.Exit(1)
	}
	slog.DebugContext(root.Context(), "finished building project for cli")

	root.AddGroup(&cobra.Group{
		ID:    "about",
		Title: "About",
	})
	root.AddCommand(versionCmd)

	return proj
}

func boostrapBuildProject(ctx context.Context) error {
	if cfgFile == "" {
		return fmt.Errorf("no config file defined")
	}

	f, foerr := os.Open(cfgFile)
	if foerr != nil {
		return fmt.Errorf("could not access config '%s' : %s", cfgFile, foerr.Error())
	}

	yb, frerr := io.ReadAll(f)
	if frerr != nil {
		return fmt.Errorf("could not read config '%s' : %s", cfgFile, frerr.Error())
	}

	tproj, umerr := config.ConfigFromYamllBytes(yb)
	if umerr != nil {
		return fmt.Errorf("Error occurred unarshalling yaml: %s \nYAML:\b%s", umerr.Error(), yb)
	}

	proj = tproj

	if valerr := proj.Validate(ctx); valerr != nil {
		return fmt.Errorf("project validation error: %s", valerr.Error())
	}

	return nil
}
