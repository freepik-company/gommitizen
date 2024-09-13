package cmd

import (
	"gommitizen/cmd/internal/cmdinit"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	argDirectory string
	argPrefix    string

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Start a repository to use gommitizen",
		Run: func(cmd *cobra.Command, args []string) {
			slog.Debug("args",
				slog.String("directory", argDirectory),
				slog.String("prefix", argPrefix),
			)
			cmdinit.Run(argDirectory, argPrefix)
		},
	}
)

func init() {
	initCmd.Flags().StringVarP(&argDirectory, "directory", "d", "", "Select a directory to initialize")
	initCmd.Flags().StringVarP(&argPrefix, "prefix", "p", "", "Select a prefix for the version file")

	rootCmd.AddCommand(initCmd)
}
