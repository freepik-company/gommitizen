package cmd

import (
	"fmt"
	"gommitizen/internal/cmdgit"
	"gommitizen/internal/config"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

type initOpts struct {
	directory string
	prefix    string
}

func Init() *cobra.Command {
	opts := initOpts{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Start a repository to use gommitizen",
		Run: func(cmd *cobra.Command, args []string) {
			initRun(opts.directory, opts.prefix)
		},
	}

	cmd.Flags().StringVarP(&opts.directory, "directory", "d", "", "Select a directory to initialize")
	cmd.Flags().StringVarP(&opts.prefix, "prefix", "p", "", "Select a prefix for the version file")

	return cmd
}

func initRun(dirPath, prefix string) {
	nDirPath, err := config.NormalizePath(dirPath)
	if err != nil {
		slog.Error(fmt.Sprintf("normalising folders: %v", err))
		os.Exit(1)
	}

	commit, err := cmdgit.GetFirstCommit()
	if err != nil {
		slog.Error(fmt.Sprintf("first commit: %v", err))
		os.Exit(1)
	}

	config := config.NewConfigVersion(nDirPath, "0.0.0", commit, prefix)
	err = config.Save()
	if err != nil {
		slog.Error(fmt.Sprintf("config: %v", err))
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("Initializing gommitizen in %s", config.GetConfigVersionFilePath()))
}
