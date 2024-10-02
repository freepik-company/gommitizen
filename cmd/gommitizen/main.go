package main

import (
	"gommitizen/cmd/gommitizen/cmd"
	"gommitizen/internal/prettylogconsole"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

type rootOpts struct {
	argDebug bool
}

func main() {
	opts := rootOpts{}

	root := &cobra.Command{
		Use:     "gommitizen",
		Version: "0.5.2",
		Short:   "A commitizen implementation for Go with multi-project support",
		Long: `A commitizen implementation for Go with multi-project support.
It only supports the conventional commits specification: https://www.conventionalcommits.org/en/v1.0.0/
Currently it only supports the bump command, but it will support the commit command soon.`,
	}

	root.Flags().BoolVar(&opts.argDebug, "debug", false, "Enable debug")

	root.AddCommand(cmd.Init())
	root.AddCommand(cmd.Bump())

	level := slog.LevelInfo
	if opts.argDebug {
		level = slog.LevelDebug
	}

	logger := slog.New(prettylogconsole.NewHandler(&slog.HandlerOptions{
		AddSource: false,
		Level:     level,
	}))
	slog.SetDefault(logger)

	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
