package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/freepik-company/gommitizen/internal/prettylogconsole"
	"github.com/freepik-company/gommitizen/internal/version"
)

type rootOpts struct {
	argDebug bool
}

func Root() *cobra.Command {
	opts := rootOpts{}

	root := &cobra.Command{
		Use:     "gommitizen",
		Version: version.GetVersion(),
		Short:   "A commitizen implementation for Go with multi-project support",
		Long: `A commitizen implementation for Go with multi-project support.
It only supports the conventional commits specification: https://www.conventionalcommits.org/en/v1.0.0/
Currently it only supports the bump command, but it will support the commit command soon.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := slog.LevelInfo
			if opts.argDebug {
				level = slog.LevelDebug
			}

			logger := slog.New(prettylogconsole.NewHandler(&slog.HandlerOptions{
				AddSource: false,
				Level:     level,
			}))
			slog.SetDefault(logger)
		},
	}
	root.PersistentFlags().BoolVar(&opts.argDebug, "debug", false, "Enable debug")

	root.AddCommand(initCmd())
	root.AddCommand(bumpCmd())
	root.AddCommand(getCmd())

	return root
}
